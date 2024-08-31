document.addEventListener("DOMContentLoaded", function() {
    fetch("/frontend/static/images/bootstrap/icons/bootstrap-icons.svg")
        .then(response => response.text())
        .then(svg => {
            const div = document.createElement("div");
            div.style.display = "none";
            div.innerHTML = svg;
            document.body.insertBefore(div, document.body.childNodes[0]);
        });
});

let currentOpenOrderId = null;
let lastClickedRow = null; // Track the last clicked row for highlighting

function toggleOrderDetails(orderId, rowElement) {
    const sideWindow = document.getElementById('side-window');

    console.log("Order ID Clicked:", orderId); // Debugging line

    if (currentOpenOrderId === orderId && sideWindow.classList.contains('open')) {
        closeSideWindow(); // Close if the same order is clicked again
        rowElement.classList.remove('highlight'); // Remove highlight on close
        console.log("Removed highlight from row."); // Debugging line
        lastClickedRow = null; // Reset last clicked row
    } else {
        if (lastClickedRow) {
            lastClickedRow.classList.remove('highlight'); // Remove highlight from the last clicked row
            console.log("Removed highlight from last clicked row."); // Debugging line
        }
        rowElement.classList.add('highlight'); // Highlight the current row
        console.log("Added highlight to current row."); // Debugging line
        lastClickedRow = rowElement; // Store the current row as last clicked
        showOrderDetails(orderId); // Show the details in the side window
        currentOpenOrderId = orderId; // Update the current open order
    }
}

function showOrderDetails(orderId) {
    let orderDetails = '';
    if (orderId === 1) {
        orderDetails = `
                <h2>Order #1 Details</h2>
                <p>Product: Apple</p>
                <p>Quantity: 10kg</p>
                <p>Due: $50</p>
                <p>Created At: Aug 27 2024 10:25pm</p>
                <p>Status: Not Paid</p>
            `;
    } else if (orderId === 2) {
        orderDetails = `
                <h2>Order #2 Details</h2>
                <p>Product: Banana</p>
                <p>Quantity: 2kg</p>
                <p>Due: $15</p>
                <p>Created At: Aug 27 2024 10:30pm</p>
                <p>Status: PAID</p>
            `;
    }

    document.getElementById('order-details').innerHTML = orderDetails;
    document.getElementById('side-window').classList.add('open');
}

function closeSideWindow() {
    document.getElementById('side-window').classList.remove('open');
    if (lastClickedRow) {
        lastClickedRow.classList.remove('highlight'); // Remove highlight when closing the side window
        console.log("Removed highlight because side window closed."); // Debugging line
    }
    currentOpenOrderId = null; // Reset the current open order
    lastClickedRow = null; // Reset last clicked row
}

function showSection(sectionId) {
    // Hide all content sections
    document.querySelectorAll('.content-section').forEach(section => {
        section.style.display = 'none'; // Hide all sections
    });

    // Show the selected content section
    document.getElementById(sectionId).style.display = 'block'; // Show the clicked section
}