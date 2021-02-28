window.onload = function() {
    if (localStorage.length != 0) {
        loadButtonOnHomepage();
    }
}

function loadButtonOnHomepage() {
    document.body.innerHTML += '<a href="/merkliste" class="merkliste"></a>';
}