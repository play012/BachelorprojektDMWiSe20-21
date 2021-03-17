window.onscroll = function() {
    var navList = document.getElementById("navList");
    var headerRect = document.getElementById("header").getBoundingClientRect();
    navList.style.top = headerRect.bottom + "px";
}

function changeUlDisplay() {
    var navList = document.getElementById("navList");
    var navIcon = document.getElementById("navIcon");

    var headerOffsetHeight = document.getElementById("header").offsetHeight;
    navList.style.top = headerOffsetHeight + "px";

    if (navList.style.display == "none") {
        navList.style.display = "block";
        navIcon.innerHTML = "X";
    } else {
        navList.style.display = "none";
        navIcon.innerHTML = "&#x2630;";
    }
}