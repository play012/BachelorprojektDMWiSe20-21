window.onload = function() {
    document.getElementById("navList").style.top = document.getElementById("header").offsetHeight;
    if (localStorage.length != 0) {
        loadButtonOnHomepage();
    }
}

function loadButtonOnHomepage() {
    document.body.innerHTML += '<a href="/merkliste" class="merkliste"></a>';
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