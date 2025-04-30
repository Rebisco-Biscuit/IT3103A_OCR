import '../custom.css';
import { gql, useQuery, useSubscription } from '@apollo/client';
import React, { useState, useEffect } from 'react';
import { DateTime } from "luxon";

const GET_PAYMENT_HISTORY = gql`
  query GetPaymentHistory($studentId: String!) {
    listPaymentHistory(studentId: $studentId) {
      transactionId
      paymentMethod
      createdAt
      courseId
      price
    }
  }
`;

const PAYMENT_CREATED = gql`
  subscription OnPaymentCreated($studentId: String!) {
    paymentCreated(studentId: $studentId) {
      transactionId
      paymentMethod
      createdAt
      courseId
      price
    }
  }
`;

const History = () => {
  const studentId = '12345';

  const [paymentHistory, setPaymentHistory] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const rowsPerPage = 8;

  // Fetch initial payment history
  const { data: queryData, loading: queryLoading, error: queryError } = useQuery(GET_PAYMENT_HISTORY, {
    variables: { studentId },
    fetchPolicy: 'network-only',
  });

  // Subscribe to new payments
  const { data: subData } = useSubscription(PAYMENT_CREATED, {
    variables: { studentId },
  });

  // On initial query load
  useEffect(() => {
    if (queryData?.listPaymentHistory) {
      const flatHistory = queryData.listPaymentHistory.flatMap((p) =>
        p.items.map((item) => ({
          transactionId: p.transactionId,
          paymentMethod: p.paymentMethod,
          createdAt: p.createdAt,
          courseId: item.courseId,
          price: item.price,
        }))
      );

      setPaymentHistory(flatHistory.reverse()); // newest on top
    }
  }, [queryData]);

  // On new payment received
  useEffect(() => {
    if (subData?.paymentCreated) {
      const p = subData.paymentCreated;

      const newItems = p.items.map((item) => ({
        transactionId: p.transactionId,
        paymentMethod: p.paymentMethod,
        createdAt: p.createdAt,
        courseId: item.courseId,
        price: item.price,
      }));

      setPaymentHistory((prev) => [...newItems, ...prev]);
    }
  }, [subData]);

  if (queryLoading) return <p style={{ marginLeft: '20px' }}>Loading your broke ass...</p>;
  if (queryError) return <p style={{ marginLeft: '20px' }}>Failed to fetch data. Cry about it.</p>;

  const startIdx = (currentPage - 1) * rowsPerPage;
  const currentRows = paymentHistory.slice(startIdx, startIdx + rowsPerPage);

  return (
    <div className="container mt-10">
      <div className="mb-10">
        <h4 style={{ fontWeight: 500 }}>Jane Doe</h4>
        <h6>Transaction history for Jane Doe</h6>
      </div>

      <div className="nocontainer" style={{ marginTop: '20px' }}>
        <table className="table">
          <thead>
            <tr>
              <th>TXN #</th>
              <th>Method</th>
              <th className="text-center">Paid At</th>
              <th className="text-center">Course - ID</th>
              <th className="text-end">Price</th>
            </tr>
          </thead>
          <tbody>
            {currentRows.map((item, index) => (
              <tr key={index}>
                <td>{item.transactionId}</td>
                <td>{item.paymentMethod}</td>
                <td className="text-center">
                  {DateTime.fromISO(item.createdAt, { zone: "utc" }).toFormat("yyyy-MM-dd hh:mm a")}
                </td>
                <td className="text-center">{item.courseId}</td>
                <td className="text-end">₱{item.price.toFixed(2)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

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
            setCurrentPage((p) => (p * rowsPerPage < paymentHistory.length ? p + 1 : p))
          }
          disabled={currentPage * rowsPerPage >= paymentHistory.length}
          className="btn"
        >
          Next
        </button>
      </div>
    </div>
  );
};

export default History;
