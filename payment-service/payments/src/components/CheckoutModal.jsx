import '../custom.css';
import React, { useState, useEffect } from "react";
import bootstrap from 'bootstrap/dist/js/bootstrap.bundle.min.js';
import { gql, useMutation } from '@apollo/client';

const CREATE_PAYMENT = gql`
  mutation CreatePayment(
    $studentId: String!,
    $items: [PaymentItemInput!]!,
    $paymentMethod: PaymentMethod!,
    $cardHolder: String,
    $cardNumber: String,
    $expiryDate: String,
    $cvv: String,
    $phoneNumber: String,
    $pin: String
  ) {
    createPayment(
      studentId: $studentId,
      items: $items,
      paymentMethod: $paymentMethod,
      cardHolder: $cardHolder,
      cardNumber: $cardNumber,
      expiryDate: $expiryDate,
      cvv: $cvv,
      phoneNumber: $phoneNumber,
      pin: $pin
    ) {
      id
      studentId
      items {
        courseId
        price
      }
      totalAmount
      transactionId
      paymentMethod
      status
      createdAt
      errorMessage
    }
  }
`;

export default function CheckoutModal({ user, totalPrice, cartItems = [], onPaymentSuccess}) {
  const [paymentMethod, setPaymentMethod] = useState("EWallet");
  const [errorMessage, setErrorMessage] = useState("");
  const [createPayment, { loading }] = useMutation(CREATE_PAYMENT);

  useEffect(() => {
    import('bootstrap/dist/js/bootstrap.bundle.min.js');
  }, []);

  const handlePaymentMethodChange = (method) => {
    setPaymentMethod(method);
    setErrorMessage("");
  };

  const handleConfirmPayment = async () => {
    setErrorMessage(""); // Clear previous error messages

    const studentId = user?.id;

    if (!studentId || cartItems.length === 0) {
      setErrorMessage("Missing user info or empty cart.");
      return;
    }

    console.log("Cart items:", cartItems); // Debugging log to check cartItems structure

    const items = cartItems.map((course) => {
      const price = course.price ? parseFloat(course.price) : 0;

      if (isNaN(price)) {
        console.warn(`Invalid price for courseId ${course.courseId}:`, course.price);
      }

      return {
        courseId: course.courseId,
        price: isNaN(price) ? 0 : price,
      };
    });

    console.log("Mapped items:", items);

    let variables = {
      studentId,
      items,
      paymentMethod,
    };

    if (paymentMethod === "EWallet") {
      const phoneNumber = document.getElementById("phoneNumber")?.value;
      const pin = document.getElementById("pin")?.value;

      if (!phoneNumber || !pin) {
        setErrorMessage("Please fill in phone number and PIN.");
        return;
      }

      variables = { ...variables, phoneNumber, pin };
    }

    if (paymentMethod === "Card") {
      const cardNumber = document.getElementById("cardHolder")?.value;
      const cardHolder = document.getElementById("cardNumber")?.value;
      const expiryDate = document.getElementById("expiryDate")?.value;
      const cvv = document.getElementById("cvv")?.value;

      if (!cardNumber || !cardHolder || !expiryDate || !cvv) {
        setErrorMessage("Complete all card fields before you pay.");
        return;
      }

      variables = { ...variables, cardHolder, cardNumber, expiryDate, cvv };
    }

    console.log("Variables being sent to mutation:", variables); // Debugging log

    try {
      const response = await createPayment({
        variables,
      });

      console.log('===response', response)
      const result = response.data.createPayment;

      if (result.status === "completed") {
        //alert("Payment successful!");

        // Close the modal programmatically
        const closeButton = document.querySelector(".btn-close[data-bs-dismiss='modal']");
        if (closeButton) {
          closeButton.click();
        }
    
        const successModal = document.createElement("div");
        successModal.className = "modal fade";
        successModal.id = "successModal";
        successModal.tabIndex = -1;
        successModal.setAttribute("aria-labelledby", "successModalLabel");
        successModal.setAttribute("aria-hidden", "true");
        successModal.innerHTML = `
          <div class="modal-dialog">
            <div class="modal-content">
              <div class="modal-header">
                <h5 class="modal-title" id="successModalLabel">Payment Successful</h5>
                <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
              </div>
              <div class="modal-body">
                Your payment has been successfully processed.
              </div>
              <div class="modal-footer">
                <button type="button" class="btn btn-primary" data-bs-dismiss="modal">Close</button>
              </div>
            </div>
          </div>
        `;

        document.body.appendChild(successModal);
        const modalInstance = new bootstrap.Modal(successModal);
        modalInstance.show();
        
        // Notify the parent component to clear the cart
        onPaymentSuccess();
      } else {
        setErrorMessage(result.errorMessage || "Payment failed. Please try again.");
      }
    } catch (err) {
      console.error("Error during payment:", err);
      setErrorMessage("An error occurred while processing the payment.");
    }
  };

  return (
    <div className="modal fade" id="checkoutModal" tabIndex="-1" aria-labelledby="exampleModalLabel" aria-hidden="true">
      <div className="modal-dialog">
        <div className="modal-content">
          <div className="modal-header" style={{ borderBottom: "none" }}>
            <h1 className="modal-title fs-5" id="exampleModalLabel">Order Details</h1>
            <button type="button" className="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>

          <div className="modal-body">
            <p>Please confirm your information is correct.</p>
            <h6>Name: <b>{user?.name || "Unknown"}</b></h6><hr />
            <h6>ID: <b>{user?.id || "N/A"}</b></h6><hr />

            <div style={{ display: "inline-flex", alignItems:"baseline" }}>
              <h6>Pay with: </h6>
              <div className="btn-group" role="group" style={{ marginLeft: "15px" }}>
                <input
                  type="radio"
                  className="btn-check"
                  name="btnradio"
                  id="btnradio1"
                  autoComplete="off"
                  defaultChecked
                  onClick={() => handlePaymentMethodChange("EWallet")}
                />
                <label className="btn" htmlFor="btnradio1">E-Wallet</label>

                <input
                  type="radio"
                  className="btn-check"
                  name="btnradio"
                  id="btnradio2"
                  autoComplete="off"
                  onClick={() => handlePaymentMethodChange("Card")}
                />
                <label className="btn" htmlFor="btnradio2">Card</label>
              </div>
            </div>

            {errorMessage && (
              <p style={{ fontSize: "12px", marginTop: "16px", color: "red" }}>
                {errorMessage}
              </p>
            )}

            {paymentMethod === "EWallet" && (
              <div id="ewalletFields" style={{ marginTop: "15px", display: "flex", justifyContent: "space-between", gap: "5px" }}>
                <input type="text" className="form-control" placeholder="Phone Number" id="phoneNumber" required />
                <input type="text" className="form-control" placeholder="4 Digit PIN" id="pin" required />
              </div>
            )}

            {paymentMethod === "Card" && (
              <div id="cardFields" style={{ marginTop: "15px" }}>
                <input type="text" className="form-control" placeholder="Card Number" id="cardHolder" style={{ marginBottom: "10px" }} required />
                <input type="text" className="form-control" placeholder="Name on card" id="cardNumber" required />
                <div className="cvv&expiryInputs" style={{ display: "flex", justifyContent: "space-between", marginTop: "10px", gap: "5px" }}>
                  <input type="text" className="form-control" placeholder="MM/YY" id="expiryDate" required />
                  <input type="text" className="form-control" placeholder="CVV" id="cvv" required />
                </div>
              </div>
            )}
          </div>

          <div className="modal-footer" style={{ justifyContent: "space-between" }}>
            <p className="text-start">
              Total Payment<br />
              <b style={{ fontSize: "20px", fontWeight: "700", color: "#7B3538" }}>₱ {totalPrice.toFixed(2)}</b>
            </p>
            <button type="button" className="btn" onClick={handleConfirmPayment} disabled={loading}>
              {loading ? "Processing..." : "Confirm and pay"}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
