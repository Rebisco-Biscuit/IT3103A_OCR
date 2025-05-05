package mq

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"payment-mod/payment-service/graph/model"
	"payment-mod/payment-service/graph/resolver"

	"github.com/go-stomp/stomp"
)

type CartEvent struct {
	Event      string  `json:"event"`
	StudentID  string  `json:"studentId"`
	CourseID   string  `json:"courseId"`
	CourseName string  `json:"courseName"`
	Price      float64 `json:"price"`
}

func StartCartConsumer(db *sql.DB) {
	for {
		err := connectAndConsume(db)
		if err != nil {
			log.Printf("Consumer error: %v. Retrying in 2 seconds...", err)
			time.Sleep(2 * time.Second)
		}
	}
}

func connectAndConsume(db *sql.DB) error {
	log.Println("Attempting STOMP connection...")

	conn, err := stomp.Dial("tcp", "localhost:61613")
	if err != nil {
		return err
	}
	defer conn.Disconnect()

	sub, err := conn.Subscribe("/queue/cart_events", stomp.AckAuto)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	log.Println("Connected and subscribed to /queue/cart_events")

	for {
		msg := <-sub.C
		if msg == nil {
			log.Println("Nil message received. Connection lost?")
			return nil // triggers outer retry loop
		}

		var event CartEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("Invalid event: %v", err)
			continue
		}

		switch event.Event {
		case "add.cart":
			handleAddToCart(db, event)
		case "remove.cart":
			handleRemoveFromCart(db, event)
		default:
			log.Printf("Unknown event type: %s", event.Event)
		}
	}
}

func handleAddToCart(db *sql.DB, event CartEvent) {
	_, err := db.Exec(`
		INSERT INTO cart (student_id, course_id, course_name, price)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT DO NOTHING`,
		event.StudentID, event.CourseID, event.CourseName, event.Price)

	if err != nil {
		log.Printf("DB error (add.cart): %v", err)
		return
	}

	log.Printf("Added course %s (₱%.2f) to cart of student %s",
		event.CourseName, event.Price, event.StudentID)

	rows, err := db.Query(`
		SELECT id, student_id, course_id, price, course_name, added_at
		FROM cart
		WHERE student_id = $1`, event.StudentID)
	if err != nil {
		log.Printf("Failed to fetch updated cart for student %s: %v", event.StudentID, err)
		return
	}
	defer rows.Close()

	var updatedCart []*model.CartItem
	for rows.Next() {
		var item model.CartItem
		err := rows.Scan(&item.ID, &item.StudentID, &item.CourseID, &item.Price, &item.CourseName, &item.AddedAt)
		if err != nil {
			log.Printf("Failed to scan cart item: %v", err)
			return
		}
		updatedCart = append(updatedCart, &item)
	}

	resolver.Mu.Lock()
	if ch, ok := resolver.CartUpdatedChannels[event.StudentID]; ok {
		ch <- updatedCart
	}
	resolver.Mu.Unlock()
}

func handleRemoveFromCart(db *sql.DB, event CartEvent) {
	_, err := db.Exec(`
		DELETE FROM cart
		WHERE student_id = $1 AND course_id = $2`,
		event.StudentID, event.CourseID)

	if err != nil {
		log.Printf("DB error (remove.cart): %v", err)
	} else {
		log.Printf("Removed course %s from cart of student %s",
			event.CourseID, event.StudentID)
	}
}
