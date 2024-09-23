// Initialize the modal with a default quantity and price input group
function addQuantityPrice() {
    const quantitiesContainer = document.getElementById("quantitiesContainer");
    const newQuantityPrice = document.createElement("div");
    newQuantityPrice.className = "input-group mb-3 quantity-price";
    newQuantityPrice.innerHTML = `
        <input type="number" class="form-control quantity" placeholder="Quantity">
        <input type="number" class="form-control price" placeholder="Price">
        <button class="btn btn-outline-secondary" onclick="removeQuantityPrice(this)">Remove</button>
    `;
    quantitiesContainer.appendChild(newQuantityPrice);
}

// Function to remove a quantity-price pair
function removeQuantityPrice(button) {
    button.parentElement.remove();
}

// Function to add the product to the products section and close the modal
async function addProduct() {
    // Fetch values from input fields
    const city = document.getElementById("cityInput").value;
    const product = document.getElementById("productInput").value;
    const quantities = document.querySelectorAll(".quantity");
    const prices = document.querySelectorAll(".price");

    // Check if required fields are filled
    if (!city || !product) {
        alert("City and Product fields are required.");
        return;
    }

    // Create data payload
    const data = {
        city: city,
        product: product,
        quantities: quantities,
        prices: prices
    };


    // Create an array to hold quantity and price pairs
    const quantityPricePairs = Array.from(quantities).map((qtyInput, index) => {
        return {
            quantity: parseFloat(qtyInput.value),
            price: parseFloat(prices[index].value)
        };
    });

    try {
        // Send POST request to the server to add a product
        const response = await fetch('/api/admin/products/add-product', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                city: city,
                product: product,
                quantityPricePairs: quantityPricePairs
            })
        });

        // Check if the response is okay
        if (!response.ok) {
            throw new Error('Failed to add product');
        }

        // Parse the JSON response
        const data = await response.json();
        console.log('Success:', data);

        // Reset the input fields in the modal
        resetModalFields();

        // Optionally, reload or update the product list on the page to include the new product
        //loadProducts();

        // Close the modal after adding the product
        const addProductModal = bootstrap.Modal.getInstance(document.getElementById('addProductModal'));
        addProductModal.hide();

    } catch (error) {
        console.error('Error:', error);
        alert('Failed to add product. Please try again.');
    }
}

// Function to remove a product section
function removeProduct(button) {
    button.closest('.product-item').remove();
}

// Function to reset the modal input fields
function resetModalFields() {
    // Reset City and Product input fields
    document.getElementById("cityInput").value = '';
    document.getElementById("productInput").value = '';

    // Reset Quantities and Prices fields
    const quantityInputs = document.querySelectorAll(".quantity-price .quantity");
    const priceInputs = document.querySelectorAll(".quantity-price .price");

    // Loop through each input and reset its value
    quantityInputs.forEach(input => input.value = '');
    priceInputs.forEach(input => input.value = '');

    // Optionally, remove additional quantity-price rows if you only want to keep one
    const quantityPriceContainer = document.getElementById("quantitiesContainer");
    const additionalRows = quantityPriceContainer.querySelectorAll(".quantity-price");

    // Keep only the first row, remove others
    additionalRows.forEach((row, index) => {
        if (index > 0) {
            row.remove();
        }
    });
}