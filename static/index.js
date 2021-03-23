window.onload = function() {
    var mapL = document.getElementById("mapL");
    var imgMapL = document.getElementById("imageMapLarge");
    var imgMapS = document.getElementById("imageMapSmall");
    var headerOffsetHeight = document.getElementById("indexHeader").offsetHeight;
    var navList = document.getElementById("navList");

    console.log(headerOffsetHeight);

    if(typeof(navList) != 'undefined' && navList != null) {
        navList.style.top = headerOffsetHeight + "px";
    }

    if(typeof(mapL) != 'undefined' && mapL != null && typeof(imgMapL) != 'undefined' && imgMapL != null && typeof(imgMapS) != 'undefined' && imgMapS != null) {
        mapL.style.maxHeight = "calc(100vh - " + headerOffsetHeight + "px)";
        imgMapL.style.top = headerOffsetHeight + "px";
        imgMapS.style.top = headerOffsetHeight + "px";
    }

    checkButtonOnHomepage();
}

window.onresize = function() {
    var mapL = document.getElementById("mapL");
    var imgMapL = document.getElementById("imageMapLarge");
    var imgMapS = document.getElementById("imageMapSmall");
    var headerOffsetHeight = document.getElementById("indexHeader").offsetHeight;
    var navList = document.getElementById("navList");

    if(typeof(navList) != 'undefined' && navList != null) {
        navList.style.top = headerOffsetHeight + "px";
    }

    if(typeof(mapL) != 'undefined' && mapL != null && typeof(imgMapL) != 'undefined' && imgMapL != null && typeof(imgMapS) != 'undefined' && imgMapS != null) {
        mapL.style.maxHeight = "calc(100vh - " + headerOffsetHeight + "px)";
        imgMapL.style.top = headerOffsetHeight + "px";
        imgMapS.style.top = headerOffsetHeight + "px";
    }
}

function checkButtonOnHomepage() {
    if (localStorage.length != 0) {
        document.getElementById("indexBody").innerHTML += '<a href="/merkliste" class="merkliste" id="merkLink"></a>';
        document.getElementById("merkLink").style.bottom = "20px";
    }
}

function changeUlDisplayIndex() {
    var navList = document.getElementById("navList");
    var navIcon = document.getElementById("navIcon");
    
    var headerOffsetHeight = document.getElementById("indexHeader").offsetHeight;
    navList.style.top = headerOffsetHeight + "px";

    if (navList.style.display == "none") {
        navList.style.display = "block";
        navIcon.innerHTML = "X";
    } else {
        navList.style.display = "none";
        navIcon.innerHTML = "&#x2630;";
    }
}