function checkFilters() {
    var checkClothing = document.getElementById("filterClothing");
    var checkTech = document.getElementById("filterTech");
    var checkFood = document.getElementById("filterFood");
    var divFood = document.getElementsByClassName("itemEssen");
    var divTech = document.getElementsByClassName("itemTechnik");
    var divClothing = document.getElementsByClassName("itemKleidung");

    if (checkClothing.checked == true) {
        for (var i=0; i<divClothing.length; i+=1) {
            divClothing[i].style.display = "none";
        }
    } else {
        for (var i=0; i<divClothing.length; i+=1) {
            divClothing[i].style.display = "flex";
        }
    }
    
    if (checkFood.checked == true) {
        for (var i=0; i<divFood.length; i+=1) {
            divFood[i].style.display = "none";
        }
    } else {
        for (var i=0; i<divFood.length; i+=1) {
            divFood[i].style.display = "flex";
        }
    }

    if (checkTech.checked == true) {
        for (var i=0; i<divTech.length; i+=1) {
            divTech[i].style.display = "none";
        }
    } else {
        for (var i=0; i<divTech.length; i+=1) {
            divTech[i].style.display = "flex";
        }
    }
}