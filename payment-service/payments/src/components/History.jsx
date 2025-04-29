import '../custom.css';
import { gql, useQuery } from '@apollo/client';
import React, { useState } from 'react';
import { DateTime } from "luxon";

const LIST_PAYMENT_HISTORY = gql`
query ListPaymentHistory($studentId: String!) {
  listPaymentHistory(studentId: $studentId) {
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

  const { loading, error, data } = useQuery(LIST_PAYMENT_HISTORY, {
    variables: { studentId },
    fetchPolicy: "network-only",
  });

  const [currentPage, setCurrentPage] = useState(1);
  const rowsPerPage = 8;

  if (loading) return <p style={{marginLeft: '20px'}}>Loading your broke ass...</p>;
  if (error) return <p style={{marginLeft: '20px'}}>Failed to fetch data. Cry about it.</p>;

  const history = data?.listPaymentHistory || [];
  const startIdx = (currentPage - 1) * rowsPerPage;
  const currentRows = history.slice(startIdx, startIdx + rowsPerPage);

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
                {DateTime.fromISO(item.createdAt, { zone: "utc" })
                    .toFormat("yyyy-MM-dd hh:mm a")}
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
        <button type="button"
          onClick={() =>
            setCurrentPage((p) =>
              p * rowsPerPage < history.length ? p + 1 : p
            )
          }
          disabled={currentPage * rowsPerPage >= history.length}
          className="btn"
        >
          Next
        </button>
      </div>
    </div>
  );
};

export default History;
