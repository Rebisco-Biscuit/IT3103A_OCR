// src/components/CartTable.jsx
import '../custom.css';
import React, { useState } from 'react';
import CheckoutModal from './CheckoutModal';
import EmptyCart from './EmptyCart'; // <- nukes the cart after successful payment

export default function CartTable({ courses = [], totalPrice = 0, onPaymentSuccess, onClearCart, onRemoveItem }) {
  const [isCartEmpty, setIsCartEmpty] = useState(false); // track if cart is empty
  const [currentPage, setCurrentPage] = useState(1); // pagination
  const rowsPerPage = 8;

  const handlePaymentSuccess = () => {
    setIsCartEmpty(true); // Clear the cart
    if (onPaymentSuccess) {
      onPaymentSuccess();
    }
  };

  const startIdx = (currentPage - 1) * rowsPerPage;
  const currentRows = courses.slice(startIdx, startIdx + rowsPerPage);

  return (
    <div className="nocontainer" style={{ marginTop: '20px' }}>
      {isCartEmpty ? (
        <EmptyCart /> // Render EmptyCart if the cart is empty
      ) : (
        <>
          {/* Clear Cart Button */}
          <div style={{ marginBottom: '20px', textAlign: 'right' }}>
            <button
              className="btn"
              onClick={onClearCart}
            >
              Clear Cart
            </button>
          </div>

          {/* Cart Table */}
          <table className="table">
            <thead>
              <tr>
                <th>Course ID</th>
                <th>Course Description</th>
                <th className="text-end">Price</th>
              </tr>
            </thead>
            <tbody>
              {currentRows.map((course, index) => (
                <tr key={index}>
                  <td>
                    <span
                      style={{
                        color: '#7B3538',
                        cursor: 'pointer',
                        fontWeight: 700,
                        marginRight: '10px',
                      }}
                      onClick={() => onRemoveItem(course.studentId, course.courseId)}
                    >
                      x
                    </span>
                    {course.courseId}
                  </td>
                  <td>{course.courseName}</td>
                  <td className="text-end">₱{course.price}</td>
                </tr>
              ))}
            </tbody>
          </table>

          <h5 style={{ fontWeight: 600, color: '#7B3538' }} className="text-end">
            TOTAL ₱{totalPrice.toFixed(2)}
          </h5>

          {/* Pagination Controls */}
          <div className="centerbtn">
            <button
              onClick={() => setCurrentPage((p) => Math.max(p - 1, 1))}
              disabled={currentPage === 1}
              className="btn"
            >
              Previous
            </button>
            <span>Page {currentPage}</span>
            <button
              onClick={() =>
                setCurrentPage((p) => (p * rowsPerPage < courses.length ? p + 1 : p))
              }
              disabled={currentPage * rowsPerPage >= courses.length}
              className="btn"
            >
              Next
            </button>
          </div>

          {/* Proceed to Checkout Button */}
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
              user={{ name: 'Jane Doe', id: '12345' }}
              totalPrice={totalPrice}
              cartItems={courses}
              onPaymentSuccess={handlePaymentSuccess}
            />
          </div>
        </>
      )}
    </div>
  );
}
