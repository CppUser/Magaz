
// Set up the SSE connection for real-time updates
let eventSource = new EventSource("/api/empl/orders");
let reconnectAttempts = 0;
const maxReconnectAttempts = 5;
let currentOpenData = null;
let lastClickedRow = null;


// function createSSEConnection() {
//     let eventSource = new EventSource("/api/empl/orders");
//
//     eventSource.onmessage = function(event) {
//         let newOrder = JSON.parse(event.data);
//
//         let existingRow = document.querySelector(`tr[data-order-id="${newOrder.id}"]`);
//         if (existingRow) {
//             existingRow.innerHTML = `
//             <td>${newOrder.id}</td>
//             <td>${newOrder.city_name}</td>
//             <td>${newOrder.product_name}</td>
//             <td>${newOrder.quantity}</td>
//             <td>${newOrder.due}</td>
//             <td>${newOrder.user_view.username}</td>
//             <td>${new Date(newOrder.created_at).toLocaleString()}</td>
//         `;
//         } else {
//             appendNewOrder(newOrder);
//         }
//
//         if (newOrder.released) {
//             removeOrderFromTable(newOrder.id);
//         }
//
//         if (currentOpenData && currentOpenData.id === newOrder.id) {
//             showOrderDetails(newOrder); // Update the side window with new details
//         }
//     };
//
//     eventSource.onerror = function(event) {
//         console.error("SSE connection error:", event);
//
//         // Reconnect logic
//         if (eventSource.readyState === EventSource.CLOSED && reconnectAttempts < maxReconnectAttempts) {
//             setTimeout(function() {
//                 console.log("Reconnecting SSE...");
//                 reconnectAttempts++;
//                 eventSource.close(); // Ensure the previous connection is closed
//                 eventSource = createSSEConnection(); // Create a new connection
//             }, 3000); // Retry after 3 seconds
//         }
//     };
//
//     return eventSource;
// }


eventSource.onmessage = function(event) {
    let newOrder = JSON.parse(event.data);

    if (newOrder.released)
    {
        removeOrderFromTable(newOrder.id);

    } else {
        let existingRow = document.querySelector(`tr[data-order-id="${newOrder.id}"]`);
        if (existingRow) {
            existingRow.innerHTML = `
            <td>${newOrder.id}</td>
            <td>${newOrder.city_name}</td>
            <td>${newOrder.product_name}</td>
            <td>${newOrder.quantity}</td>
            <td>${newOrder.due}</td>
            <td>${newOrder.user_view.username}</td>
            <td>${new Date(newOrder.created_at).toLocaleString()}</td>
        `;
        } else {
            appendNewOrder(newOrder);
        }


        if (currentOpenData && currentOpenData.id === newOrder.id) {
            showOrderDetails(newOrder); // Update the side window with new details
        }
    }




    // if (newOrder.released) {
    //     removeOrderFromTable(newOrder.id);
    // } else {
    //     appendNewOrder(newOrder);
    // }
};

eventSource.onerror = function(event) {
    console.error("SSE connection error:", event);

    // Reconnect logic
    if (eventSource.readyState === EventSource.CLOSED && reconnectAttempts < maxReconnectAttempts) {
        setTimeout(function() {
            console.log("Reconnecting SSE...");
            reconnectAttempts++;
            eventSource.close(); // Ensure the previous connection is closed
            eventSource = new EventSource("/api/empl/orders"); // Create a new connection
        }, 3000); // Retry after 3 seconds
    }
};


function appendNewOrder(order) {
    let tableBody = document.getElementById("orders-tbody");
    console.log(order);

    // Check if order already exists in the table
    let existingRow = document.querySelector(`tr[data-order-id="${order.id}"]`);
    if (existingRow) {
        console.log(`Order #${order.id} already exists in the table.`);
        return; // Prevent duplicate entries
    }

    // Create a new row for the new order
    let row = document.createElement("tr");
    row.setAttribute("data-order-id", order.id);

    row.innerHTML = `
        <td>${order.id}</td>
        <td>${order.city_name}</td>
        <td>${order.product_name}</td>
        <td>${order.quantity}</td>
        <td>${order.due}</td>
        <td>${order.user_view.username}</td>
        <td>${new Date(order.created_at).toLocaleString()}</td>
    `;

    row.addEventListener("click", function() {
        console.log(`Order #${order.id} clicked`);
        toggleOrderDetails(order, row); // Pass the order object and the row element
    });

    // Append the new row to the table
    tableBody.appendChild(row);

}


function toggleOrderDetails(order, rowElement) {
    const sideWindow = document.getElementById('side-window');

    // Check if the same order is clicked again to close the window
    if (currentOpenData === order && sideWindow.classList.contains('open')) {
        closeSideWindow(); // Close if the same order is clicked again
        rowElement.classList.remove('highlight'); // Remove highlight on close
        lastClickedRow = null; // Reset last clicked row
    } else {
        // Close the side window and remove highlight from the last row if needed
        if (lastClickedRow) {
            lastClickedRow.classList.remove('highlight'); // Remove highlight from the last clicked row
        }

        // Highlight the current row and show the order details in the side window
        rowElement.classList.add('highlight'); // Highlight the current row
        lastClickedRow = rowElement; // Store the current row as last clicked
        showOrderDetails(order); // Show the details in the side window
        currentOpenData = order; // Update the current open order
    }
}

function showOrderDetails(order) {
    console.log(`showOrderDetails called for order #${order.id}`);
    document.getElementById('order-details').innerHTML = `

            <!-- Modal for displaying available addresses -->
            <div id="assignAddressModal" class="modal-overlay table-responsive">
                <div class="modal-header">
                     <h2>Выберите адрес</h2>
                </div>
                 <div class="modal-body">
                <table class="table-secondary address-table modal-body">
                <thead>
                    <tr>
                        <th scope="col">#</th>
                        <th scope="col">Город</th>
                        <th scope="col">Продукт</th>
                        <th scope="col">Кол</th>
<!--                        <th scope="col">Работник</th>-->
                        <th scope="col">Фото</th>
                        <th scope="col">Добавлен</th>
                    </tr>
                </thead>
                 <tbody class="table-group-divider">
                 
                 </tbody>
                </table>
                </div>
                <div class="modal-footer">
                    <button id="assign-btn" class="btn btn-success" disabled>Присвоить адрес</button>
                </div>
            </div>
            
            
        <h2>Заказ #${order.id}</h2>
        <p>Город: ${order.city_name}</p>
        <p>Продукт: ${order.product_name}</p>
        <p>Количество: ${order.quantity}</p>
        <p>К оплате: ${order.due}</p>
        <p data-section="created_at">Добавлен: ${new Date(order.created_at).toLocaleString()}</p>
  
        <!-- Client Information -->
        <p class="collapsible-toggle" data-section="client" onclick="toggleDetails('client-details')">Клиент: ${order.user_view.username}</p>
        <div id="client-details" class="collapsible" style="display: none;">
            <p>User ID: ${order.user_view.id}</p>
            <p>Чат ID: ${order.user_view.chat_id}</p>
        </div>
       
        <!-- Payment Information -->
        <p class="collapsible-toggle" data-section="payment" onclick="toggleDetails('payment-info')">Метод оплаты: ${order.payment_method.payment_category}</p>
        <div id="payment-info" class="collapsible" style="display: none;">
            <p>Банк: ${order.payment_method.card_payment.bank_name}</p>
        </div>
       
        <!-- Address Information -->
        <p class="collapsible-toggle" data-section="address" onclick="toggleDetails('address-info')">
            Присвоеный адресс: #${order.address && order.address.id ? order.address.id : 'Не присвоен'}
        </p>
        <div id="address-info" class="collapsible" style="display: none;">
            ${
        order.address && order.address.id
            ? `
                <p>Описание: ${order.address.description}</p>
                <p>Фото: ${order.address.image}</p>
                <p>Количество: ${order.address.quantity}</p>
                <p>Дата: ${new Date(order.address.added_at).toLocaleString()}</p>
                <p>Добавлен: имя работника</p>
                                
               
                <!-- Edit button if address is assigned -->
                <button onclick="editAddress(${order.address.id})" class="btn btn-secondary">Редактировать адрес</button>
                <button onclick="reassignAddress(${order.id})" class="btn btn-warning">Присвоить новый</button>
                `
            : `
                <!-- Assign button if no address is assigned -->
                
                <button id="assign-btn" class="btn btn-success" onclick="assignAddress(${order.id})">Присвоить адрес</button>
                
                
               
                
                `
    }
        </div>
        
        <!-- Release and Decline Buttons -->
        <div class="action-buttons">
            <button onclick="showReleasePopup(${order.id})" class="btn btn-success">Подтвердить</button>
            <button onclick="showDeclinePopup(${order.id})" class="btn btn-danger">Отклонить</button>
        </div>
        
        
    </div>
</div>
        
    `;

    // Show the side window
    document.getElementById('side-window').classList.add('open');
}

function toggleDetails(id) {
    const section = document.getElementById(id);
    const toggle = document.querySelector(`[onclick="toggleDetails('${id}')"]`);

    // Toggle the display of the section
    if (section.style.display === "none" || section.style.display === "") {
        section.style.display = "block";
        toggle.classList.add("open"); // Add the "open" class to rotate the arrow
    } else {
        section.style.display = "none";
        toggle.classList.remove("open"); // Remove the "open" class to rotate the arrow back
    }
}


function closeSideWindow() {
    document.getElementById('side-window').classList.remove('open');
    if (lastClickedRow) {
        lastClickedRow.classList.remove('highlight');
        lastClickedRow = null;
    }
}

function showSection(sectionId) {
    // Hide all content sections
    document.querySelectorAll('.content-section').forEach(section => {
        section.style.display = 'none'; // Hide all sections
    });

    // Show the selected content section
    document.getElementById(sectionId).style.display = 'block'; // Show the clicked section
}

function assignAddress(orderId) {
    document.getElementById('assignAddressModal').style.display = 'block';

    fetch(`/api/empl/orders/address?orderId=${orderId}`)
        .then(response => response.json())
        .then(data => {
            console.log(data)
            const addressTableBody = document.querySelector('#assignAddressModal tbody'); // Select the table body
            addressTableBody.innerHTML = '';
            data.addresses.forEach(address => {
                console.log(data.addresses)
                const row = document.createElement('tr');

                const idCell = document.createElement('td');
                idCell.textContent = address.id;

                const cityCell = document.createElement('td');
                cityCell.textContent = address.city;

                const productCell = document.createElement('td');
                productCell.textContent = address.product;

                //TODO:
                const quantityCell = document.createElement('td');
                quantityCell.textContent = address.quantity;

                // const employeeCell = document.createElement('td');
                // employeeCell.textContent = address.added_by

                //TODO: add thumbnail image
                const imageCell = document.createElement('td');
                const thumbnail = document.createElement('img');
                thumbnail.src = `/api/get/images/${encodeURIComponent(address.image)}`;
                //thumbnail.alt = 'Thumbnail';
                thumbnail.style.width = '50px'; // Thumbnail size
                thumbnail.style.height = '50px';
                thumbnail.style.cursor = 'pointer';
                thumbnail.onclick = () => openFullImage(`/api/get/images/${encodeURIComponent(address.image)}`); // Open full image
                imageCell.appendChild(thumbnail);

                const addedAtCell = document.createElement('td');
                addedAtCell.textContent = new Date(address.added_at).toLocaleString(); // Format date

                row.appendChild(idCell);
                row.appendChild(cityCell);
                row.appendChild(productCell);
                row.appendChild(quantityCell);
                // row.appendChild(employeeCell);
                row.appendChild(imageCell);
                row.appendChild(addedAtCell);

                row.style.cursor = 'pointer';
                row.onclick = () => {
                    // Remove the green highlight from all rows
                    const rows = document.querySelectorAll('.table-group-divider tr');
                    rows.forEach(r => r.style.backgroundColor = ''); // Reset all rows

                    // Highlight the clicked row
                    row.style.backgroundColor = 'lightgreen'; // Highlight in green

                    // Enable the assign button and assign the address
                    document.getElementById('assign-btn').disabled = false;
                    document.getElementById('assign-btn').onclick = () => assignSelectedAddress(orderId, address.id);
                };

                addressTableBody.appendChild(row);
            });
            document.getElementById('assignAddressModal').style.display = 'block'; // Show the modal
        })
        .catch(error => console.error('Error fetching addresses:', error));
}

function assignSelectedAddress(orderId, addressId) {
    // Send a POST request to the backend to assign the address
    fetch(`/api/empl/orders/address/assign`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            order_id: orderId,
            address_id: addressId,
        }),
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                document.getElementById('assignAddressModal').style.display = 'none';

                const existingRow = document.querySelector(`tr[data-order-id="${orderId}"]`);
                if (existingRow) {
                    existingRow.remove();
                }

                showOrderDetails(data.order);

               // updateOrderTable(data.order);

            } else {
                alert('Failed to assign address: ' + data.error);
            }
        })
        .catch(error => {
            console.error('Error assigning address:', error);
        });
}

function openFullImage(imageUrl) {
    const modal = document.createElement('div');
    modal.style.position = 'fixed';
    modal.style.top = '0';
    modal.style.left = '0';
    modal.style.width = '100%';
    modal.style.height = '100%';
    modal.style.backgroundColor = 'rgba(0, 0, 0, 0.8)';
    modal.style.display = 'flex';
    modal.style.justifyContent = 'center';
    modal.style.alignItems = 'center';
    modal.style.cursor = 'pointer';
    modal.style.zIndex = '1300';

    const fullImage = document.createElement('img');
    fullImage.src = imageUrl;
    fullImage.style.maxWidth = '80%';
    fullImage.style.maxHeight = '80%';

    modal.appendChild(fullImage);
    modal.onclick = () => modal.remove(); // Close on click

    document.body.appendChild(modal);
}

function editAddress(addressId) {
    // Logic to edit an assigned address (e.g., open a modal or form)
    console.log("Editing address with ID:", addressId);
}

function reassignAddress(orderId) {
    // Logic to reassign the address (e.g., open a modal or form to choose a new address)
    console.log("Reassigning address for order ID:", orderId);
}

function closeModal() {
    document.getElementById('assignAddressModal').style.display = 'none';
}

function showReleasePopup(orderId) {
    const releaseForm = `
        <div class="release-popup">
            <h3>Release Order</h3>
            <p>Please confirm the order has been paid.</p>
            <button onclick="confirmRelease(${orderId})" class="btn btn-success">Confirm</button>
            <button onclick="cancelRelease()" class="btn btn-secondary">Cancel</button>
        </div>
    `;

    // Show the release form in a modal or popup
    document.body.insertAdjacentHTML('beforeend', releaseForm);
}

function confirmRelease(orderId) {
    // Send release request to the backend
    fetch(`/api/empl/orders/release/${orderId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            //TODO:'Authorization': 'Bearer ' + token // Add your authentication token here
        }
    })
        .then(response => response.json())
        .then(data => {
            if (data.message === "Order released successfully") {
                closeSideWindow();
                removeOrderFromTable(orderId); // Remove the order from the table after releasing
            } else {
                alert("Failed to release order: " + data.error);
            }
        })
        .catch(err => {
            console.error("Error releasing order:", err);
        });

    // Remove the release popup
    cancelRelease();
}

function cancelRelease() {
    const popup = document.querySelector('.release-popup');
    if (popup) {
        popup.remove();
    }
}



function showDeclinePopup(orderId) {
    const declineForm = `
        <div class="decline-popup">
            <h3>Decline Order</h3>
            <p>Select a reason for declining the order:</p>
            <select id="decline-reason">
                <option value="Not paid">Not paid</option>
                <option value="Run out of time">Run out of time</option>
                <option value="Out of stock">Out of stock</option>
                <option value="Other">Other</option>
            </select>
            <button onclick="submitDecline(${orderId})" class="btn btn-danger">Submit Decline</button>
            <button onclick="cancelDecline()" class="btn btn-secondary">Cancel</button>
        </div>
    `;

    // Show the decline form in a modal or popup
    document.body.insertAdjacentHTML('beforeend', declineForm);
}

function submitDecline(orderId) {
    const reason = document.getElementById('decline-reason').value;
    console.log("Order declined with reason:", reason);

    // Logic to handle declining the order
    declineOrder(orderId, reason);

    // Remove the decline popup
    cancelDecline();
}

function cancelDecline() {
    const popup = document.querySelector('.decline-popup');
    if (popup) {
        popup.remove();
    }
}

function declineOrder(orderId, reason) {
    // Logic to decline the order with the selected reason
    console.log(`Order ID: ${orderId}, Reason: ${reason}`);
}

function removeOrderFromTable(orderId) {
    let tableBody = document.getElementById("orders-tbody");
    let rows = tableBody.getElementsByTagName("tr");

    for (let i = 0; i < rows.length; i++) {
        const row = rows[i];
        const rowOrderId = row.firstElementChild.textContent.trim();

        // Convert both to strings for comparison
        if (rowOrderId === orderId.toString()) {
            tableBody.removeChild(row); // Remove the row from the table
            console.log(`Order #${orderId} removed from the table.`);
            break;
        }
    }
}

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