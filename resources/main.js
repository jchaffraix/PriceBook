'use strict';

class Price {
  constructor(price, quantity, unit) {
    // Convert all input to numbers.
    // This makes equals works as we get numbers
    // from our model but stings from the UI.
    this.price = Number(price);
    this.quantity = Number(quantity);
    this.unit = unit;
  }

  equals(otherPrice) {
    return this.price === otherPrice.price
      && this.quantity === otherPrice.quantity
      && this.unit === otherPrice.unit;
  }

  get pricePerQuantity() {
    return Math.round(100 * this.price / this.quantity) / 100;
  }

  get pricePerQuantityText() {
    // TODO: Stop hardcoding '$'.
    return "$" + this.pricePerQuantity + "/" + this.unit;
  }
}

class Product {
  constructor(name, brand, price, place, date) {
    this.name = name;
    this.brand = brand;
    this.price = price;
    this.place = place;
    this.date = date;
  }
}

function generateDateInputCell(date) {
  return generateSingleInputCell(date, "date");
}

function generateSingleTextInputCell(name) {
  return generateSingleInputCell(name, "text");
}

function generateSingleInputCell(value, type) {
  var td = document.createElement("td");
  var semanticInput = document.createElement("div");
  semanticInput.classList.add("ui");
  semanticInput.classList.add("transparent");
  semanticInput.classList.add("input");
  var input = document.createElement("input");
  input.type = type;
  input.value = value;
  input.addEventListener("change", updateModel);
  semanticInput.appendChild(input);
  td.appendChild(semanticInput);
  return td;
}

function createUnitSelection(unit) {
  var unitNode = document.createElement("select");
  unitNode.classList.add("ui");
  unitNode.classList.add("dropdown");

  var item = document.createElement("option");
  item.value = "item";
  item.appendChild(document.createTextNode("item"));
  unitNode.appendChild(item);

  var pound = document.createElement("option");
  pound.value = "pound";
  pound.appendChild(document.createTextNode("pound"));
  unitNode.appendChild(pound);

  var oz = document.createElement("option");
  oz.value = "oz";
  oz.appendChild(document.createTextNode("oz"));
  unitNode.appendChild(oz);

  unitNode.value = unit;
  unitNode.addEventListener("change", updateModel);
  return unitNode;
}

function generatePriceAndQuantityCell(price) {
  var td = document.createElement("td");
  var semanticInput = document.createElement("div");
  semanticInput.classList.add("ui");
  semanticInput.classList.add("transparent");
  semanticInput.classList.add("input");
  var currency = document.createElement("span");
  // TODO: Remove this hard-coding.
  currency.appendChild(document.createTextNode("$"));
  semanticInput.appendChild(currency);

  var priceNode = document.createElement("input");
  priceNode.type = "number";
  priceNode.classList.add("price");
  priceNode.pattern = "/^\d+\.?\d*$/";
  priceNode.value = price.price;
  priceNode.addEventListener("keypress", function(e) { if(e.target.value.length==5) e.preventDefault(); });
  priceNode.addEventListener("change", updateModel);
  semanticInput.appendChild(priceNode);

  var forSeparator = document.createElement("span");
  forSeparator.classList.add("forSeparator");
  forSeparator.appendChild(document.createTextNode("for"));
  semanticInput.appendChild(forSeparator);

  var quantity = document.createElement("input");
  quantity.type = "number";
  quantity.value = price.quantity;
  quantity.classList.add("quantity");
  quantity.pattern = "\d+";
  quantity.addEventListener("keypress", function(e) { if(e.target.value.length==3) e.preventDefault(); });
  quantity.addEventListener("change", updateModel);
  semanticInput.appendChild(quantity);

  td.appendChild(semanticInput);
  td.appendChild(createUnitSelection(price.unit));
  return td;
}

function findEnclosingRow(node) {
  // Find the enclosing row
  var row = node.parentNode;
  while (row.tagName !== "TR") {
    row = row.parentNode;
  }
  return row;
}

function updateModel(e) {
  const row = findEnclosingRow(e.target);
  var product = row.product;

  // Always update the name, brand and the place.
  const name = row.childNodes[1].getElementsByTagName("input")[0].value;
  const brand = row.childNodes[2].getElementsByTagName("input")[0].value;
  const place = row.childNodes[5].getElementsByTagName("input")[0].value;
  product.name = name;
  product.brand = brand;
  product.place = place;

  const priceCell = row.childNodes[3];
  const price = priceCell.getElementsByTagName("input")[0].value;
  const quantity = priceCell.getElementsByTagName("input")[1].value;
  const unit = priceCell.getElementsByTagName("select")[0].value;
  const newPrice = new Price(price, quantity, unit);
  if (!product.price.equals(newPrice)) {
    product.price = newPrice;

    var pricePerQuantityCell = row.childNodes[4];
    var newPricePerQuantityCell = generatePricePerQuantityCell(newPrice);
    row.replaceChild(newPricePerQuantityCell, pricePerQuantityCell);
  }

  // Update the date unless it was updated by the user.
  const dateCell = row.childNodes[6];
  const date = dateCell.getElementsByTagName("input")[0].value;
  var shouldUpdateDate = product.date === date;
  if (shouldUpdateDate) {
    product.date = formatToday();
    const newDateCell = generateDateInputCell(product.date);
    row.replaceChild(newDateCell, dateCell);
  } else {
    product.date = date;
  }
  saveProductList();
}

function generatePricePerQuantityCell(price) {
  var td = document.createElement("td");
  td.appendChild(document.createTextNode(price.pricePerQuantityText));
  return td;
}

function formatToday() {
  // toISOString returns the time, which we won't be accepted
  // by <input type="date"> so we just isolate the date part.
  return new Date().toISOString().split("T")[0];
}

function generateIconsCell() {
  var cell = document.createElement("td");
  var removeIcon = document.createElement("i");
  removeIcon.classList.add("icon");
  removeIcon.classList.add("trash");
  removeIcon.classList.add("alternate");
  removeIcon.classList.add("outline");
  removeIcon.addEventListener("touch", removeProduct);
  removeIcon.addEventListener("click", removeProduct);
  cell.appendChild(removeIcon);

  var duplicateIcon = document.createElement("i");
  duplicateIcon.classList.add("icon");
  duplicateIcon.classList.add("copy");
  duplicateIcon.classList.add("outline");
  duplicateIcon.addEventListener("touch", duplicateProduct);
  duplicateIcon.addEventListener("click", duplicateProduct);
  cell.appendChild(duplicateIcon);
  return cell;
}

function createTableRow(product) {
  // Save the model in the row so we can update it.
  var row = document.createElement("tr");
  row.product = product;
  row.appendChild(generateIconsCell());
  var nameCell = generateSingleTextInputCell(product.name);
  row.appendChild(nameCell);
  row.appendChild(generateSingleTextInputCell(product.brand));
  row.appendChild(generatePriceAndQuantityCell(product.price));
  row.appendChild(generatePricePerQuantityCell(product.price));
  row.appendChild(generateSingleTextInputCell(product.place));
  row.appendChild(generateDateInputCell(product.date));
  return row;
}

// Returns the products that match the search
// query. If it empty, we return everything.
function productsFilteredBySearch() {
  const filter = document.getElementById("search").value;
  var filteredProducts = Array.from(productList);
  if (filter !== "") {
    filteredProducts = filteredProducts.filter(product => product.name.startsWith(filter));
  }
  return filteredProducts.sort(function(productA, productB) {
    if (productA.name < productB.name)
      return -1;
    if (productA.name > productB.name)
      return 1;
    return 0;
  });
}

function populateTable() {
  const tbody = document.getElementsByTagName("tbody")[0];
  var new_tbody = document.createElement("tbody");
  const products = productsFilteredBySearch();
  for (var product of products) {
    new_tbody.appendChild(createTableRow(product));
  }
  tbody.parentNode.replaceChild(new_tbody, tbody);
}

function addNewProduct() {
  const product = new Product("", "", new Price(0, 1, ""), "", formatToday());
  productList.push(product);
  // TODO: We need to migrate the element to a better position once it is populated.
  const tbody = document.getElementsByTagName("tbody")[0];
  const row = createTableRow(product);
  tbody.insertBefore(row, tbody.firstChild);
  row.childNodes[1].getElementsByTagName("input")[0].focus();
  saveProductList();
}

function removeProduct(e) {
  const row = findEnclosingRow(e.target);
  const productToRemove = row.product;
  for (var i = 0; i < productList.length; ++i) {
    const product = productList[i];
    if (product === productToRemove) {
      productList.splice(i, 1);
      break;
    }
  }
  saveProductList();
  populateTable();
}

function duplicateProduct(e) {
  const row = findEnclosingRow(e.target);
  const productToDuplicate = row.product;
  for (var i = 0; i < productList.length; ++i) {
    const product = productList[i];
    if (product === productToDuplicate) {
      productList.splice(i, 0, productToDuplicate);
      break;
    }
  }
  saveProductList();
  populateTable();
}

var productList = new Array();

function parseProductList(products) {
  for (const product of products) {
    const price = product.price;
    productList.push(new Product(product.name, product.brand, new Price(price.price, price.quantity, price.unit), product.place, product.date));
  }
}

function populateProductList() {
  var products = window.localStorage.getItem("products");
  // For first time users, localStorage will return null.
  if (products === null)
    return;

  parseProductList(JSON.parse(products));
}

function saveProductList() {
  window.localStorage.setItem("products", JSON.stringify(productList));
}

function initializeApp() {
  populateProductList();
  populateTable();
  $("#add").click(addNewProduct);
  // This listens to input to react while the user is typing.
  document.getElementById("search").addEventListener("input", populateTable);
}

window.addEventListener("load", initializeApp);
