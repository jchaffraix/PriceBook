const dbName = "price_book_db";
const dbVersion = 1;
const testProducts = [
  { "product": "Pear", "brand": "Barlett", "unit": "Pound", "number": 1, "price": 1.89, "tag": "organic", "shop": "Lucky's", "date": 3/2/2019 },
  { "product": "Hummus", "brand": "Sabra", "unit": "Oz", "number": 17, "price": 6.99, "tag": "", "shop": "Lucky's", "date": 3/2/2019 },
];

var request = indexedDB.open(dbName, dbVersion);

request.onerror = function(event) {
  // Handle errors.
};
request.onsuccess = function(event) {
  db = event.target.result;
};
request.onupgradeneeded = function(event) {
  var db = event.target.result;

  // Create an objectStore to hold information about our customers. We're
  // going to use "ssn" as our key path because it's guaranteed to be
  // unique - or at least that's what I was told during the kickoff meeting.
  var objectStore = db.createObjectStore("products", { keyPath: "id", autoIncrement:true });

  // Create an index to search customers by name. We may have duplicates
  // so we can't use a unique index.
  objectStore.createIndex("brand", "branch", { unique: false });

  // TODO: Index per store?

  // Use transaction oncomplete to make sure the objectStore creation is 
  // finished before adding data into it.
  objectStore.transaction.oncomplete = function(event) {
    // Store values in the newly created objectStore.
    var productObjectStore = db.transaction("products", "readwrite").objectStore("products");
    testProducts.forEach(function(product) {
      productObjectStore.add(product);
    });
  };
};
