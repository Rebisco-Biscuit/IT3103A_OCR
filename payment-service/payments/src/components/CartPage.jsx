// src/components/CartPage.jsx
import '../custom.css';
import React, { useState } from 'react';
import CartTable from './CartTable';
import EmptyCart from './EmptyCart'; // <- nukes the cart after successful payment

export default function CartPage() {
  const [cartItems, setCartItems] = useState([
    { id: 'Course008', description: 'Otto', price: '9999.00' },
  ]);

  // Calculate total price dynamically
  const totalPrice = cartItems.reduce(
    (total, item) => total + parseFloat(item.price.replace(/[^\d.-]/g, '')),
    0
  );

  // Handle payment success by emptying the cart
  const handlePaymentSuccess = () => {
    setCartItems([]); // Empty the cart
  };

  return (
    <div className="container mt-10">
        <div className="mb-10">
          <h4 style={{ fontWeight: 500 }}>Jane Doe</h4>
          <h6>Cart of listed courses for Jane Doe</h6>
      </div>

      {/* Conditionally render CartTable or EmptyCart */}
      {cartItems.length === 0 ? (
        <EmptyCart /> // Show EmptyCart if the cart is empty
      ) : (
        <CartTable
          courses={cartItems}
          totalPrice={totalPrice}
          onPaymentSuccess={handlePaymentSuccess} // Pass the success handler to CartTable
        />
      )}
    </div>
  );
}
