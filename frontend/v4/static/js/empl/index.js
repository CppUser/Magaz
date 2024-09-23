const sideLinks = document.querySelectorAll('.container .sidebar .side-menu li a:not(.logout)');

sideLinks.forEach(item => {
    const li = item.parentElement;
    item.addEventListener('click', () => {
        sideLinks.forEach(i => {
            i.parentElement.classList.remove('active');
        })
        li.classList.add('active');
    })
});

const menuBar = document.querySelector('.sidebar .bx.bx-menu');
const sideBar = document.querySelector('.sidebar');

menuBar.addEventListener('click', () => {
    sideBar.classList.toggle('close');
});


window.addEventListener('resize', () => {
    if (window.innerWidth < 768) {
        sideBar.classList.add('close');
    } else {
        sideBar.classList.remove('close');
    }
});

const toggler = document.getElementById('theme-toggle');

toggler.addEventListener('change', function () {
    if (this.checked) {
        document.body.classList.add('dark');
    } else {
        document.body.classList.remove('dark');
    }
});

function showSection(sectionId) {
    // Hide all content sections
    document.querySelectorAll('.content main .content-section').forEach(section => {
        section.style.display = 'none'; // Hide all sections
    });

    // Show the selected content section
    document.getElementById(sectionId).style.display = 'block'; // Show the clicked section
}

let currentOpenData = null;
let lastClickedRow = null;

function toggleRowDetails(order, rowElement) {
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
        orderDetails(order); // Show the details in the side window
        currentOpenData = order; // Update the current open order
    }
}

function orderDetails(order) {
    document.getElementById('order-details').innerHTML = `

            <!-- Modal for displaying available addresses -->
            <div id="assignAddressModal" class="modal-overlay">
            <button class="close-btn" onclick="closeAddressWindow()">
                    <i class="bx bxs-x-square"></i>
                </button>
                <div class="modal-header">
                     <h2>Выберите адрес</h2>
                </div>
                 <div class="modal-body">
                <table class="address-table modal-body">
                <thead>
                    <tr>
                        <th scope="col">Кол</th>
                        <th scope="col">Фото</th>
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
        <p class="collapsible-toggle" data-section="client" onclick="showDetails('client-details')">Клиент: ${order.user_view.username}</p>
        <div id="client-details" class="collapsible" style="display: none;">
            <p>User ID: ${order.user_view.id}</p>
            <p>Чат ID: ${order.user_view.chat_id}</p>
        </div>
       
        <!-- Payment Information -->
        <p class="collapsible-toggle" data-section="payment" onclick="showDetails('payment-info')">Метод оплаты: ${order.payment_method.payment_category}</p>
        <div id="payment-info" class="collapsible" style="display: none;">
            <p>Банк: ${order.payment_method.card_payment.bank_name}</p>
            <p>Ссылка: <a href="https://${order.payment_method.card_payment.bank_url}" target="_blank" rel="noopener noreferrer">${order.payment_method.card_payment.bank_url}</a></p>
            <p>Карта #: ${order.payment_method.card_payment.card_number}</p>
            <p>ФИО: ${order.payment_method.card_payment.first_name}  ${order.payment_method.card_payment.last_name}</p>
            <p>СБП: ${order.payment_method.card_payment.quick_pay}</p>
        </div>
       
        <!-- Address Information -->
        <p class="collapsible-toggle" data-section="address" onclick="showDetails('address-info')">
            Присвоеный адресс: #${order.address && order.address.id ? order.address.id : 'Не присвоен'}
        </p>
        <div id="address-info" class="collapsible" style="display: none;">
            ${
        order.address && order.address.id
            ? `
                <p>Описание: ${order.address.description}</p>
                <p>Фото: 
                <img src="/api/get/images/${order.address.image}" 
                     alt="Thumbnail" 
                     style="width: 100px; height: auto; cursor: pointer;" 
                     onclick="openImage(\`/api/get/images/${encodeURIComponent(order.address.image)}\`);">
                </p>
                <p>Количество: ${order.address.quantity}</p>
                <p>Дата: ${new Date(order.address.added_at).toLocaleString()}</p>
                <p>Добавлен: имя работника</p>
                                
               
                <!-- Edit button if address is assigned -->
                <!--<button onclick="editAddress(${order.address.id})" class="btn btn-secondary">Редактировать адрес</button>
                <button onclick="reassignAddress(${order.id})" class="btn btn-warning">Присвоить новый</button>-->
                `
            : `
                <!-- Assign button if no address is assigned -->
                
                <button id="assign-btn" class="btn btn-success" onclick="bindNewAddress(${order.id})">Присвоить адрес</button>
                
                
               
                
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

function closeSideWindow() {
    const sideWindow = document.getElementById('side-window');
    if (sideWindow.classList.contains('open')) {
        sideWindow.classList.remove('open');
    }

    // Ensure that the last clicked row is de-highlighted
    if (lastClickedRow) {
        lastClickedRow.classList.remove('highlight');
        lastClickedRow = null;
    }
}

function closeAddressWindow() {
    document.getElementById('assignAddressModal').style.display = 'none';
}

function showDetails(sec) {
    const section = document.getElementById(sec);
    const toggle = document.querySelector(`[onclick="showDetails('${sec}')"]`);

    if (section.style.display === "none" || section.style.display === "") {
        section.style.display = "block";
        toggle.classList.add("open"); // Add the "open" class to rotate the arrow
    } else {
        section.style.display = "none";
        toggle.classList.remove("open"); // Remove the "open" class to rotate the arrow back
    }
}


function bindNewAddress(orderId) {
    document.getElementById('assignAddressModal').style.display = 'block';

    fetch(`/api/empl/orders/address?orderId=${orderId}`)
        .then(response => response.json())
        .then(data => {
            const addressTableBody = document.querySelector('#assignAddressModal tbody'); // Select the table body
            addressTableBody.innerHTML = '';
            data.addresses.forEach(address => {
                const row = document.createElement('tr');


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
                thumbnail.onclick = () => openImage(`/api/get/images/${encodeURIComponent(address.image)}`); // Open full image
                imageCell.appendChild(thumbnail);

                row.appendChild(quantityCell);
                // row.appendChild(employeeCell);
                row.appendChild(imageCell);


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

//TODO Should send message via web socket
function assignSelectedAddress(orderId, addressId) {

    if (window.isWebSocketOpen && isWebSocketOpen()) {
        sendEvent('update_address',{"orderId": orderId, "addressId": addressId})
    }else {
        alert('Срединение прервано , обновите страницу для отоброжения результата');
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



}

// Function to show the release popup
function showReleasePopup(orderId) {
    const modal = document.createElement('div');
    modal.classList.add('release-modal'); // Use a unique CSS class for the release popup
    modal.innerHTML = `
        <div class="popup-content">
            <h3>Подтвердить заказ #${orderId}</h3>
            <button class="btn btn-success" onclick="confirmRelease(${orderId})">Подтвердить</button>
            <button class="btn btn-danger" onclick="closeReleasePopup()">Отмена</button>
        </div>`;
    document.body.appendChild(modal);
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
                closeReleasePopup();
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
    const modal = document.createElement('div');
    modal.classList.add('decline-modal'); // Use a unique CSS class for the decline popup
    modal.innerHTML = `
        <div class="popup-content">
            <h3>Отклонить заказ #${orderId}</h3>
             <label for="declineReason">Причина:</label>
            <select id="declineReason">
                <option value="not_paid">Не оплачено</option>
                <option value="customer_cancelled">Клиент отменил</option>
                <option value="out_of_stock">Нет в наличии</option>
            </select>
            <br>
            <button class="btn btn-danger" onclick="confirmDecline(${orderId})">Отклонить</button>
            <button class="btn btn-secondary" onclick="closeDeclinePopup()">Отмена</button>
        </div>`;
    document.body.appendChild(modal);
}

function confirmDecline(orderId) {
    const declineReason = document.getElementById("declineReason").value; // Get selected reason

    // Send decline request to the backend
    fetch(`/api/empl/orders/decline/${orderId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            //TODO: 'Authorization': 'Bearer ' + token // Add your authentication token here
        },
        body: JSON.stringify({
            reason: declineReason
        })
    })
        .then(response => response.json())
        .then(data => {
            if (data.message === "Order declined successfully") {
                closeSideWindow();
                removeOrderFromTable(orderId); // Remove the order from the table after declining
                closeDeclinePopup(); // Close the decline popup
            } else {
                alert("Failed to decline order: " + data.error);
            }
        })
        .catch(err => {
            console.error("Error declining order:", err);
        });
}

// Function to close the release popup
function closeReleasePopup() {
    const modal = document.querySelector('.release-modal'); // Select the release modal
    if (modal) {
        modal.remove(); // Remove the modal from the DOM
    }
}

function closeDeclinePopup() {
    const modal = document.querySelector('.decline-modal');
    if (modal) {
        modal.remove(); // Remove the modal from the DOM
    }
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
        const rowOrderId = row.dataset.orderId; // Assuming you've added a data attribute for order ID

        // Ensure orderId is compared as a string or number
        if (rowOrderId && rowOrderId.toString() === orderId.toString()) {
            tableBody.removeChild(row); // Remove the row from the table
            console.log(`Order #${orderId} removed from the table.`);
            break;
        }
    }
}

function openImage(imageUrl) {
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