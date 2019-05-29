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
  constructor(key, name, brand, price, place, date) {
    // The key is an internal identifier used to deduplicate
    // objects in the DB.
    this._key = key;
    this.name = name;
    this.brand = brand;
    this.price = price;
    this.place = place;
    this.date = date;
    // TODO: Add an LMT to handle conflicting writes.
  }

  static creatEmpty() {
    return new Product(this._uniqueKey(), "", "", new Price(0, 1, ""), "", formatToday());
  }

  static createFromJSON(json) {
    const price = new Price(json.price.price, json.price.quantity, json.price.unit); 
    return new Product(json._key, json.name, json.brand, price, json.place, json.date);
  }

  clone() {
    // Clone generates a new key to ensure the objects are saved as different objects.
    return new Product(this._uniqueKey(), this.name, this.brand, this.price, this.place, this.date);
  }

  static _uniqueKey() {
    return Math.random().toString(36).substring(2)
                 + (new Date()).getTime().toString(36);
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
  const filter = new RegExp(document.getElementById("search").value, "i");
  var filteredProducts = Array.from(productList);
  if (filter !== "") {
    filteredProducts = filteredProducts.filter(product => product.name.search(filter) != -1);
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
  const product = Product.creatEmpty();
  productList.push(product);
  // TODO: We need to migrate the element to a better position once it is populated.
  // TODO: There is a general issue around renaming and sorting where elements whose
  // name changed are not moved to their appropriate position until we regenerate
  // the table. This also means that |productList| is not sorted, which could have
  // performance issues.
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
      productList.splice(i, 0, productToDuplicate.clone());
      break;
    }
  }
  saveProductList();
  populateTable();
}

var productList = new Array();

function parseProductList(products) {
  var newProductList = new Array();
  for (const product of products) {
    newProductList.push(Product.createFromJSON(product));
  }
  return newProductList;
}

function fetchProductList() {
  // TODO: I should keep track of the LMT to avoid
  // cloberring local file.
  $("#loader").addClass("active");
  $.ajax({
    dataType: "json",
    url: "/store",
    data: null,
    success: function (products) {
      // TODO: Add an animation?
      $("#loader").removeClass("active");
      $("#success").removeClass("disabled");
      $("#failureHolder").addClass("disabled");
      console.log("Fetch success");
      // TODO: Sanity checks?
      productList = parseProductList(products);
      populateTable();
    },
    error: function (jXHR, textStatus, errorThrown) {
      // TODO: Add an animation?
      $("#loader").removeClass("active");
      $("#success").addClass("disabled");
      console.log("Fetch failed");
      console.log(jXHR.status);
      if (jXHR.status === 401) {
        $("#login").removeClass("disabled");
      } else {
        $("#failureHolder").removeClass("disabled");
      }
    },
  });
}

function populateProductList() {
  // Load from localStorage first.
  var products = window.localStorage.getItem("products");
  // For first time users, localStorage will return null.
  if (products !== null) {
    productList = parseProductList(JSON.parse(products));
  }

  fetchProductList();
}

window.updateTimer = null;
function saveProductList() {
  window.localStorage.setItem("products", JSON.stringify(productList));
  // Kickstart the 10s update timer.
  // TODO: This is not useful if the user is not logged in.
  if (window.updateTimer === null) {
    window.updateTimer = window.setTimeout(function() {
      // TODO: Share this code.
      $("#loader").addClass("active");
      $("#success").addClass("disabled");
      $("#failureHolder").addClass("disabled");
      $.ajax({
        method: "POST",
        url: "/store",
        data: JSON.stringify(productList),
        success: function (products) {
          // TODO: Add an animation?
          $("#loader").removeClass("active");
          $("#success").removeClass("disabled");
          $("#failureHolder").addClass("disabled");
          window.updateTimer = null;
          console.log("Saved success");
        },
        error: function (jXHR, textStatus, errorThrown) {
          // TODO: Add an animation?
          $("#loader").removeClass("active");
          $("#success").addClass("disabled");
          window.updateTimer = null;
          console.log("Failed");
          console.log(jXHR.textStatus);
          console.log(jXHR.errorThrown);
          if (jXHR.status === 401) {
            $("#login").removeClass("disabled");
          } else {
            $("#failureHolder").removeClass("disabled");
          }
        },
      });
    }, 5000 /* 5 seconds */);
  }
}

function initializeApp() {
  populateProductList();
  populateTable();
  $("#add").click(addNewProduct);
  // TODO: This will force the user to always click on the button.
  // I can do better than this by querying directly and getting a 401.
  $("#login").click(function() {
    var win = window.open("/login", "login_window");
    win.onunload = function() {
      $("#login").css("display", "none");
      fetchProductList();
    };
  });
  // This listens to input to react while the user is typing.
  document.getElementById("search").addEventListener("input", populateTable);
}

window.addEventListener("load", initializeApp);
window.onbeforeunload = function() {
  if (window.updateTimer !== null) {
    return "Pending changes";
  }
}
