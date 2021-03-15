var eventHandler = document.getElementsByClassName("addItemBtn");
var iBox = document.getElementById('infoBox');




var objects = [iBox,iBox,iBox,iBox,iBox,iBox,iBox,iBox,iBox];
var counters = [0,0,0,0,0,0,0,0,0,0,0,,0,0];
var transformObj;
var elem;

var indexClick = 0;
var o = 1;


eventHandler.onclick = function() {

    //infoBox clonen, neu benennen
    elem = iBox.cloneNode(true);
    elem.id = "infoBox" + indexClick.toString();
    iBox.appendChild(elem); // An Ibox wird elem angehängt
    
    //img erstellen und anhängen
    var img = document.createElement('img');
    img.id = "img";
    img.src = "addBasket.png";
    elem.appendChild(img);  // An elem wird img gehängt
    //img stylen
    img.style.width = "100px";
    img.style.height = "100px";
    img.style.marginTop = "50px";
    
    //elem in Array setzen -> keine überschneidungen
    objects[indexClick] = elem;
    var obj = objects[indexClick];
    
    //stepCount in Array setzen -> keine überschneidungen
    counters[indexClick] = indexClick;
    var stepCount = counters[indexClick];
    
    console.log("Element erstellt: %s",elem.id);
    console.log("IndexClick: %d", indexClick);
    
    //Style Object
    obj.style.position = "relative";
    obj.style.width = "100px";
    obj.style.height = "100px";
    obj.style.textAlign = "center";

    //Set parameters = 0    
    transformObj = 0;
    //Count -> 10
    if (indexClick >= 10) indexClick = 0;
     else indexClick++;
    

    // -> Moving Object (through moveObj function)
    var myInterval = setInterval(moveObj, 12);  




    function moveObj(){
        //up to stepCount -> move Object
        if (stepCount <= 140) {
            transformObj = 0;
            transformObj += stepCount;
            console.log("Object %d, TransformValue: %d",indexClick ,transformObj)
            // MovementLength/HoldLength
            if (stepCount < 50) {
                obj.style.transform = "translateX(-"+transformObj.toString()+"px)";
            }
            //step Size (=animation speed)
            stepCount+=4;

            // Use last steps for OpacityFade
        } else {
            obj.style.opacity = o.toString();
            o -= 0.05;
            
            if (o <= 0) {
                o = 1;
                stepCount = 0;
                console.log("Object %d wird geloescht!", indexClick);
                indexClick = 0;
                obj.style.display = "none";
                clearInterval(myInterval);
            }
        }
    
    }

}

    //sleep function 
    function sleep(milliseconds) {
        const date = Date.now();
        let currentDate = null;
        do {
          currentDate = Date.now();
        } while (currentDate - date < milliseconds);
    }


