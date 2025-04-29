// src/components/CartTable.jsx
import '../custom.css';
import React, { useState } from 'react';
import CheckoutModal from './CheckoutModal';
import EmptyCart from './EmptyCart'; // <- nukes the cart after successful payment

export default function CartTable({ courses = [], totalPrice = 0, onPaymentSuccess }) {
  const [isCartEmpty, setIsCartEmpty] = useState(false); // State to track if the cart is empty

  const handlePaymentSuccess = () => {
    setIsCartEmpty(true); // Set the cart to empty after a successful payment
    if (onPaymentSuccess) {
      onPaymentSuccess();
    }
  };

  return (
    <div className="nocontainer" style={{ marginTop: '20px' }}>
      {isCartEmpty ? (
        <EmptyCart /> // Render EmptyCart if the cart is empty
      ) : (
        <>
          <table className="table">
            <thead>
              <tr>
                <th>Course ID</th>
                <th>Course Description</th>
                <th className="text-end">Price</th>
              </tr>
            </thead>
            <tbody>
              {courses.map((course, index) => (
                <tr key={index}>
                  <td>x {course.id}</td>
                  <td>{course.description}</td>
                  <td className="text-end">₱{course.price}</td>
                </tr>
              ))}
            </tbody>
          </table>

          <h5 style={{ fontWeight: 600 }} className="text-end">
            TOTAL ₱{totalPrice.toFixed(2)}
          </h5>

          <div
            style={{
              marginTop: '25px',
              display: 'flex',
              justifyContent: 'center',
              alignItems: 'center',
            }}
          >
            <button
              type="button"
              className="btn"
              data-bs-toggle="modal"
              data-bs-target="#checkoutModal"
            >
              Proceed to checkout
            </button>

            <CheckoutModal
              user={{ name: "Jane Doe", id: "12345" }}
              totalPrice={totalPrice}
              cartItems={courses}
              onPaymentSuccess={handlePaymentSuccess} // Pass the success handler to CheckoutModal
            />
          </div>
        </>
      )}
    </div>
  );
}
