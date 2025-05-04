// src/components/CartPage.jsx
import '../custom.css';
import React, { useState, useEffect } from 'react';
import { gql, useQuery, useSubscription, useMutation } from '@apollo/client';
import CartTable from './CartTable';
import EmptyCart from './EmptyCart';

const GET_CART = gql`
  query GetCart($studentId: String!) {
    getCart(studentId: $studentId) {
      id
      studentId
      courseId
      price
      courseName
      addedAt
    }
  }
`;

const CART_UPDATED = gql`
  subscription CartUpdated($studentId: String!) {
    cartUpdated(studentId: $studentId) {
      id
      studentId
      courseId
      price
      courseName
      addedAt
    }
  }
`;

const CLEAR_CART = gql`
  mutation ClearCart($studentId: String!) {
    clearCart(studentId: $studentId)
  }
`;

export default function CartPage() {
  const studentId = '12345'; // will be replaced with the actual student ID
  const [cartItems, setCartItems] = useState([]);

  // fetch initial cart items from course service
  const { data: queryData, loading: queryLoading, error: queryError } = useQuery(GET_CART, {
    variables: { studentId },
    fetchPolicy: 'network-only',
  });

  const { data: subData } = useSubscription(CART_UPDATED, {
    variables: { studentId },
  });

  const [clearCart] = useMutation(CLEAR_CART, {
    variables: { studentId },
    onCompleted: () => {
      console.log('Cart cleared successfully after payment');
    },
    onError: (error) => {
      console.error('Failed to clear cart:', error);
    },
  });

  useEffect(() => {
    if (queryData?.getCart) {
      setCartItems(queryData.getCart);
    }
  }, [queryData]);

  useEffect(() => {
    if (subData?.cartUpdated) {
      setCartItems(subData.cartUpdated);
    }
  }, [subData]);

  const totalPrice = cartItems.reduce((total, item) => total + item.price, 0);

  // clear the cart after payment
  const handlePaymentSuccess = () => {
    clearCart(); // automatically clear the cart after payment
    setCartItems([]); // shows empty cart UI
  };

  if (queryLoading) return <p>Loading your cart...</p>;
  if (queryError) return <p>Failed to load cart. Please try again later.</p>;

  return (
    <div className="container mt-10">
      <div className="mb-10">
        <h4 style={{ fontWeight: 500 }}>Jane Doe</h4>
        <h6>Cart of listed courses for Jane Doe</h6>
      </div>

      {/* Conditionally render CartTable or EmptyCart */}
      {cartItems.length === 0 ? (
        <EmptyCart />
      ) : (
        <CartTable
          courses={cartItems}
          totalPrice={totalPrice}
          onPaymentSuccess={handlePaymentSuccess} // Pass the success handler to CartTable
          onClearCart={() => clearCart()} 
        />
      )}
    </div>
  );
}
