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
      courseName
      price
      addedAt
    }
  }
`;

console.log('===CART_UPDATED', CART_UPDATED)

const CART_REMOVEITEM = gql`
  mutation RemoveFromCart($studentId: String!, $courseId: String!) {
    removeFromCart(studentId: $studentId, courseId: $courseId)
  }
`;

const CLEAR_CART = gql`
  mutation ClearCart($studentId: String!) {
    clearCart(studentId: $studentId)
  }
`;

export default function CartPage() {
  const studentId = '12345'; // Replace with the actual student ID
  const [cartItems, setCartItems] = useState([]);

  // Fetch initial cart items
  const { data: queryData, loading: queryLoading, error: queryError } = useQuery(GET_CART, {
    variables: { studentId },
    fetchPolicy: 'network-only',
  });

  // Subscribe to cart updates
  const { data: subData } = useSubscription(CART_UPDATED, {
    variables: { studentId },
  });


  // Mutation to clear the cart
  const [clearCart] = useMutation(CLEAR_CART, {
    variables: { studentId },
    onCompleted: () => {
      console.log('Cart cleared successfully after payment');
      setCartItems([]); // Update the UI to reflect the empty cart
    },
    onError: (error) => {
      console.error('Failed to clear cart:', error);
    },
  });

  // Mutation to remove a specific course from the cart
  const [removeFromCart] = useMutation(CART_REMOVEITEM, {
    onCompleted: () => {
      console.log('Course removed from cart successfully');
    },
    onError: (error) => {
      console.error('Failed to remove course from cart:', error);
    },
  });

  // Update cart items when the query data is loaded
  useEffect(() => {
    if (queryData?.getCart) {
      console.log('===setCartItem?')
      setCartItems(queryData.getCart);
    }
  }, [queryData]);
  console.log('===queryData', queryData);
  console.log('===subDatam', subData)
  // Update cart items when a subscription event is received

  useEffect(() => {
    if (subData && subData.cartUpdated) {
      console.log('===Data', subData);
      setCartItems(subData.cartUpdated);
    }
  }, [subData]);
  

  // Calculate total price dynamically
  const totalPrice = cartItems.reduce((total, item) => total + item.price, 0);

  // Clear the cart after payment
  const handlePaymentSuccess = () => {
    clearCart(); // Automatically clear the cart after payment
  };

  // Remove a specific course from the cart
  const handleRemoveItem = (studentId, courseId) => {
    removeFromCart({
      variables: { studentId, courseId },
    });
    // Update the cart items locally
    setCartItems((prevItems) => prevItems.filter((item) => item.courseId !== courseId));
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
          onClearCart={() => clearCart()} // Pass the clearCart handler
          onRemoveItem={handleRemoveItem} // Pass the remove item handler
        />
      )}
    </div>
  );
}
