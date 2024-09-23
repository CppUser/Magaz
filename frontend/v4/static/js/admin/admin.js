const sideLinks = document.querySelectorAll('.sidebar .side-menu li a:not(.logout)');

sideLinks.forEach(item => {
    const li = item.parentElement;
    item.addEventListener('click', () => {
        sideLinks.forEach(i => {
            i.parentElement.classList.remove('active');
        })
        li.classList.add('active');
    })
});

function displaySection(sectionId) {
    document.querySelectorAll('.content-section > div').forEach(section => {
        section.style.display = 'none';
    });

    const selectedSection = document.getElementById(sectionId);
    if (selectedSection) {
        console.log("Displaying section:", sectionId);
        selectedSection.style.display = 'block';
    }
}

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

function removeQuantityPrice(button) {
    button.parentElement.remove();
}