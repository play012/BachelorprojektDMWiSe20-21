window.onload = function() {
    var headerOffsetHeight = document.getElementById("header").offsetHeight;
    var navListTop = document.getElementById("navList").style.top;

    if (navListTop && headerOffsetHeight != null) {
        navListTop = headerOffsetHeight;
    }

    console.log(localStorage.length);    
    checkButtonOnHomepage();
}

function checkButtonOnHomepage() {
    if (localStorage.length != 0) {
        document.getElementById("indexBody").innerHTML += '<a href="/merkliste" class="merkliste"></a>';
    }
}

function changeUlDisplay() {
    var navList = document.getElementById("navList");
    var navIcon = document.getElementById("navIcon");
    if (navList.style.display == "none") {
        navList.style.display = "block";
        navIcon.innerHTML = "X";
    } else {
        navList.style.display = "none";
        navIcon.innerHTML = "&#x2630;";
    }
}