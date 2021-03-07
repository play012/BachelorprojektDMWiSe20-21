window.onload = function() {
    if (window.location.pathname == "/merkliste") {
        checkLocalStorageForNull();
        if (localStorage.length != 0) {
            loadItemsFromLocalStorage();
            showClearAllBtn();
        }
    }
}

function addToLocalStorage(idNum, itemCategory, itemOffer, itemStore) {
    let item = [idNum, itemCategory, itemOffer, itemStore];
    let itemStr = JSON.stringify(item);

    if (eval('document.getElementById("addItemBtn' + idNum +'") != null')) {
        eval('localStorage.setItem("Item' + idNum + '", itemStr);');
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

        imgDiv.className = "imgWrapper";
        if (itemArr[j][1] == "Technik") {
            itemDiv.className = "itemTechnik";
            imgDiv.innerHTML = '<img src="/static/tp.png" id="techImg">';
        } else if (itemArr[j][1] == "Essen") {
            itemDiv.className = "itemEssen";
            imgDiv.innerHTML = '<img src="/static/tp.png" id="foodImg">';
        } else if (itemArr[j][1] == "Kleidung") {
            itemDiv.className = "itemKleidung";
            imgDiv.innerHTML = '<img src="/static/tp.png" id="clothingImg">';
        } else {
            itemDiv.style.display = "none";
        }

        textDiv.className = "textWrapper";
        var offerText = itemArr[j][2];
        var storeText = itemArr[j][3];
        textDiv.innerHTML = '<p>' + offerText + '</p><p>' + storeText + '</p>';

        btnDiv.className = "imgWrapper";
        btnDiv.innerHTML = '<button class="delItemBtn" id="delItemBtn' + itemArr[j][0] + '" onclick="removeItemsFromLocalStorage(' + itemArr[j][0] + ')"></button>';

        itemDiv.appendChild(imgDiv);
        itemDiv.appendChild(textDiv);
        itemDiv.appendChild(btnDiv);
        document.body.appendChild(itemDiv);
    }
}

function removeItemsFromLocalStorage(id) {
    eval('localStorage.removeItem("Item' + id +'");');
    document.location.reload();
}

function showClearAllBtn() {
    var clearAllBtn = document.getElementById("clearAllBtn");
    clearAllBtn.style.display = "block";
}

function clearAllItems() {
    localStorage.clear();
    document.location.reload();
}

function submitForm() {
    document.getElementById("sortForm").submit();
}