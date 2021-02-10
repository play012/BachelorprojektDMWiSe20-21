window.onload = function() {
    if (window.location.pathname == "/merkliste") {
        checkLocalStorageForNull();
        if (localStorage.length != 0) {
            loadItemsFromLocalStorage();
        }
    }
}

function addToLocalStorage(idNum, itemCategory, itemOffer, itemStore) {
    let item = [itemCategory, itemOffer, itemStore];
    let itemStr = JSON.stringify(item);
    var compArr = [];

    for (var k = 0; k < localStorage.length; k++) {
        compArr = JSON.parse(localStorage.getItem(localStorage.key(k)));
        console.log("what");
        if (compArr == item) {
            eval('document.getElementById("addItemBtn' + idNum +'").id = "delItemBtn' + idNum +'";');
            eval('document.getElementById("delItemBtn' + idNum +'").onclick = removeItemFromLocalStorage();');
            console.log("dumped item");
            return;
        }

        if (eval('document.getElementById("addItemBtn' + idNum +'") != null);')) {
            console.log("item saved successfully");
            eval('localStorage.setItem("Item' + localStorage.length + '", itemStr);');
            eval('document.getElementById("addItemBtn' + idNum +'").id = "delItemBtn' + idNum +'";');
            eval('document.getElementById("delItemBtn' + idNum +'").onclick = removeItemFromLocalStorage();');
        }
        console.log("item save failed");
    }
}

function compareItemsToLocalStorage(idNum, itemCatListed, itemOfferListed, itemStoreListed) {
    var itemsListed = [itemCatListed, itemOfferListed, itemStoreListed];
    var compItem = [];
    for (var l = 0; l < localStorage.length; l++) {
        compItem = JSON.parse(localStorage.getItem(localStorage.key(l)));
        if (compItem == itemsListed) {
            eval('document.getElementById("addItemBtn' + idNum +'").id = "delItemBtn' + idNum +'";');
            eval('document.getElementById("delItemBtn' + idNum +'").onclick = removeItemFromLocalStorage();');
        }
    }
}

function checkLocalStorageForNull() {
    var noItemsText = document.getElementById("noItemsStored");
    if (localStorage.length == 0) {
        noItemsText.style.display = "block";
        noItemsText.innerHTML = "Aktuell sind keine Angebote in der Merkliste gespeichert.";
    }
}

function loadItemsFromLocalStorage() {
    var itemArr = [];
    for (var i = 0; i < localStorage.length; i++) {
        var itemPos = localStorage.key(i);
        itemArr[i] = JSON.parse(localStorage.getItem(itemPos));
    }
    
    console.log(itemArr);

    for (var j = 0; j < itemArr.length; j++) {
        var itemDiv = document.createElement("div");
        var imgDiv = itemDiv.cloneNode(true);
        var textDiv = itemDiv.cloneNode(true);
        var btnDiv = itemDiv.cloneNode(true);
        itemDiv.className = "itemBox";

        imgDiv.className = "imgWrapper";
        if (itemArr[j][0] == "Technik") {
            imgDiv.innerHTML = '<img src="/static/tp.png" id="techImg">';
        } else if (itemArr[j][0] == "Essen") {
            imgDiv.innerHTML = '<img src="/static/tp.png" id="foodImg">';
        } else if (itemArr[j][0] == "Kleidung") {
            imgDiv.innerHTML = '<img src="/static/tp.png" id="clothingImg">';
        } else {
            itemDiv.style.display = "none";
        }

        textDiv.className = "textWrapper";
        var offerText = itemArr[j][1];
        var storeText = itemArr[j][2];
        textDiv.innerHTML = '<p>' + offerText + '</p><p>' + storeText + '</p>';

        btnDiv.className = "imgWrapper";
        btnDiv.innerHTML = '<button class="delItemBtn" id="delItemBtn' + j + '" onclick="removeItemFromLocalStorage()"></button>';

        itemDiv.appendChild(imgDiv);
        itemDiv.appendChild(textDiv);
        itemDiv.appendChild(btnDiv);
        document.getElementById("scrollBox").appendChild(itemDiv);
    }
}

function removeItemFromLocalStorage() {
    
}