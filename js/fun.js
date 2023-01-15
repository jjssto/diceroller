// This code is licensed under the MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.


export { 
    hasFocus, 
    highlightOwnRolls,
    addCol,
    formatTime,
    init,
    createRowDice,
    setResultHeaderDice,
    initCookieConsent,
};

var hasFocus;
var highlightOwnRolls;
var displayDice;

function detectFocus() {
    hasFocus = true
    window.addEventListener('focus', function (event) {
        hasFocus = true
        window.addEventListener('blur', function (event) {
            hasFocus = false
        })
    })
}

function settingVisibility() {
    var element = document.getElementById("f_setting");
    if (sessionStorage.getItem("is_owner") == "true") {
        element.style.visibility = 'visible';
    } else {
        element.style.visibility = 'hidden';
    }
}


function roomSettingForm() {
    //insertColorOptions()
    var settingsForm = document.getElementById("f_setting");
    settingsForm.addEventListener("submit", function (event) {
        event.preventDefault();
        var roomName = document.getElementById("f_setting_name").value;
        var colorSelect = document.getElementById("f_setting_color");
        var color = colorSelect.options[colorSelect.selectedIndex].value;
        var roomId = sessionStorage.getItem("room_id");
        fetch("/roomSettings", {
            method: "POST",
            headers: {
                "Accept": "application/json",
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                "room_id": roomId,
                "room_name": roomName,
                "color": color
            })
        })
        .then(res => { 
            window.location.reload()
        })
    })
}

function changeColor(color) {
    document.querySelector(".header").style.backgroundColor = color;
    document.querySelector(".footer").style.backgroundColor = color;
    document.querySelectorAll("thead").forEach(element => {
       element.style.backgroundColor = color; 
    });
    document.querySelectorAll("th").forEach(element => {
       element.style.backgroundColor = color; 
    });

    document.getElementById("f_setting_color").value = color;
}

function setColor() {
    window.addEventListener("load", function (event) {
        var color = this.sessionStorage.getItem("color");
        if (color.length > 0) {
            changeColor(color);
        }
    })
    window.addEventListener("reset", function (event) {
        var color = this.sessionStorage.getItem("color");
        if (color.length > 0) {
            changeColor(color);
        }
    })

}

function setLink() {
    var link = document.getElementById("a_link");
    const loc = window.location.href;
    link.href = loc
    link.textContent = loc
}

function addCol(row, text) {
    const col = row.insertCell()
    col.appendChild(document.createTextNode(text))
}



function insertRolls(createRow, data_raw) {
    const data = JSON.parse(data_raw)
    const tbody = document.getElementById("tbody_rolls")
    const first_row = tbody.firstChild
    var roll_id
    if (first_row) {
        roll_id = first_row.firstChild.firstChild.textContent
    } else {
        roll_id = -1
    }
    var counter = 0
    var lastRow
    for (let i in data) {
        const drow = data[i]
        sessionStorage.setItem("last_roll", i)
        if (i > roll_id) { 
           const row = createRow(drow, i)
           lastRow = tbody.appendChild(row)
        }
        counter++
    }
    if (lastRow) {
        lastRow.scrollIntoView(true)
    }
}

function getRolls(createRow) {
    if (!hasFocus) {
        return
    }
    var target = location.href.replace("room/", "rolls/");
    const last_roll = sessionStorage.getItem("last_roll");
    if (last_roll != "") {
        target += "/" + last_roll
    }
    fetch(target, {
        method: "GET",
    })
    .then(response => {
        let moreData = parseInt(response.headers.get("more_data"));
        response.text().then(data => {
                insertRolls(createRow, data);
        })
        if (moreData > 0) {
            getRolls(createRow);
        }
    });
}

function setHighlightOwnRolls() {
    const element = document.getElementById("i_highlight_own_rolls")
    if (element == null) {
        return
    }
    if (element.checked) {
        highlightOwnRolls = true;
    } else {
        highlightOwnRolls = false;
    }
}

function setDisplayDice() {
    const element = document.getElementById("i_display_dice")
    if (element == null) {
        return
    }
    if (element.checked) {
        displayDice = true;
        setResultHeaderDice();
    } else {
        displayDice = false;
    }
}

function formatTime(timestamp) {
    const d = new Date(parseInt(timestamp) * 1000) 
    const h = d.getHours()
    const m = d.getMinutes()
    const hh = ((h < 10) ? '0' : '') + h 
    const mm = ((m < 10) ? '0' : '') + m 
    return hh + ':' + mm
}

function init(createRow, rollDice) {

    window.addEventListener("DOMContentLoaded", () => {
        setHighlightOwnRolls()
        setDisplayDice()
        setColor()
        detectFocus()
        settingVisibility()
        roomSettingForm()
        setLink()
        smallScreen()

        document.getElementById("f_roll").addEventListener("submit", event => {
            event.preventDefault()
            rollDice()
        })
       
        var resetButton = document.querySelector("#b_reset")
        if (resetButton != null) {
            resetButton.addEventListener("click", (event) => {
                reset(event);
            });
        }
        initRadioButtons();
        initActionButtons(rollDice);
        initCookieConsent();
        initAllDiceForm();
    })

    window.addEventListener("load", () => {
        sessionStorage.setItem("last_roll", "");
        reset(null)
        if (displayDice) {
            getRolls(createRowDice);
            window.setInterval(() => {getRolls(createRowDice)}, 1000);
        } else {
            getRolls(createRow);
            window.setInterval(() => {getRolls(createRow)}, 1000);
        }
    });
}

function createDie(p, die) {
    var div = document.createElement("div")
    div.classList.add("icon")
    var img = document.createElement("img")
    var result 
    var eyes 
    if (die.E == "0") {
        eyes = "10"
        var r = parseInt(die.R) 
        if (r == "0") {
            result = "00"
        } else {
            result = (r * 10).toString()
        }
    } else {
        eyes = die.E
        result = die.R
    }
    img.src = "/pic/d" + eyes + ".svg"
    div.appendChild(img)
    var nbr = document.createElement("div")
    nbr.classList.add("centered")
    nbr.textContent = result
    div.appendChild(nbr)
    p.appendChild(div)
}

function createColDice(row, dice) {
    var col = document.createElement("td")
    var div = document.createElement("div")
    div.classList.add("dice_icons")
    for(var i = 0; i < dice.length; i++) {
        createDie(div, dice[i])
    }
    col.appendChild(div)
    row.appendChild(col)
}

function createRowDice(drow, id) {
    var row = document.createElement("tr")
    const dat = drow[0]
    const isOwnRoll = drow[1]
    if (isOwnRoll == "1" && highlightOwnRolls) {
        row.className = "my_roll"
    }
    addCol(row, id)
    addCol(row, formatTime(dat.T))
    const colP = row.insertCell()
    colP.appendChild(document.createTextNode(dat.P))
    colP.className = "my_name"
    addCol(row, dat.A)
    createColDice(row, dat.D)
    addCol(row, dat.R)
    return row
}

function setResultHeaderDice() {
    const thead = document.querySelector(".t_result").querySelector("thead")
    const cols = thead.querySelectorAll("th")
    for (var i = 0; i < cols.length; i++) {
        cols[i].remove()
    }
    var h1 = document.createElement("th")
    var h2 = document.createElement("th")
    var h3 = document.createElement("th")
    var h4 = document.createElement("th")
    var h5 = document.createElement("th")
    var h6 = document.createElement("th")
    h1.textContent = "#"
    h2.textContent = "Time"
    h3.textContent = "Char"
    h4.textContent = "Action"
    h5.textContent = "Dice"
    h5.style.width = "400px"
    h6.textContent = "Result"
    thead.appendChild(h1)
    thead.appendChild(h2)
    thead.appendChild(h3)
    thead.appendChild(h4)
    thead.appendChild(h5)
    thead.appendChild(h6)
}
    
    
function reset(event) {
    if (event != null) {
        event.preventDefault()
    }
    var numberSelects = document.querySelectorAll(".dd_dice_nbr")
    for (var i = 0; i < numberSelects.length; i++ ) {
        numberSelects[i].value = 0
    }
    var selectedRadios = document.querySelectorAll(".r_selected");
    for (var i = 0; i < selectedRadios.length; i++ ) {
        selectedRadios[i].classList.remove('r_selected');
    }
    var defaultRadios = document.querySelectorAll(".r1")
    for (var i = 0; i < defaultRadios.length; i++ ) {
        defaultRadios[i].classList.add('r_selected');
    }
    let mod = sessionStorage.getItem("mod");
    if (mod != null) {
        sessionStorage.setItem("mod", 0);
    }
    var action = document.querySelector("#f_action")
    if (action != null) {
        action.value = ""
    }
}

function initRadioButtons() {
    var radioButtons = document.querySelectorAll(".b_radio");
    for (var i = 0; i < radioButtons.length; i++ ) {
        radioButtons[i].addEventListener("click", (event) => {
            let selected = document.querySelectorAll('.r_selected')
            for (let j = 0; j < selected.length; j++) {
                if (selected[j].name == event.target.name ) {
                    selected[j].classList.remove("r_selected");
                    break;
                }
            }
            event.target.classList.add('r_selected')
            let element = document.getElementById(event.target.name) 
            if (element != null && !(typeof element === 'undefined')) {
                element.value = event.target.value
            } else {
                sessionStorage.setItem("mod", event.target.value);
            }
        })    
    }
    let radioElements = document.getElementsByClassName("roll_form_radio"); 
    for (var i = 0; i < radioElements.length; i++ ) {
        let idStr = radioElements[i].id.substring(6);
        let element = document.getElementById(idStr);
        if (element == null) continue;
         element.addEventListener("change", (event) => {
            let radio = document.getElementById("radio_" + event.target.id)
            for (var j = 0; j < radio.children.length; j++ ) {
                if (radio.children[j].value == event.target.selectedIndex) {
                    radio.children[j].classList.add("r_selected")
                } else {
                    radio.children[j].classList.remove("r_selected")
                }
            }
        })
    }
}

function initActionButtons(rollFunction) {
    var actionButtons = document.querySelectorAll(".button_action");
    var actionInput = document.getElementById("f_action");
    for (var i = 0; i < actionButtons.length; i++ ) {
        actionButtons[i].addEventListener("click", (event) => {
            let oldVal = actionInput.value;
            actionInput.value = event.target.textContent;;
            rollFunction();
            actionInput.value = oldVal;
        })
    }
}

function initCookieConsent() {

    // if cookie already exists => leave function
    if (checkCookie()) {
        return
    }

    // otherwise initialise cookie consent form
    document.querySelector(".cookie_consent")
        .classList.replace("hidden", "visible");
    document.getElementById("f_cookie_consent")
        .addEventListener("submit", (event) => {
            event.preventDefault();
            const checkbox = document.getElementById("i_cookie_consent")
            if (checkbox != null && checkbox.checked ) {
                if (setCookie()) {               
                    document.querySelector(".cookie_consent")
                        .classList.replace("visible", "hidden");
                }
            }
    })
    document.getElementById("i_cookie_consent")
        .addEventListener("click", (event) => {
            if (event.target.checked) {
                document.getElementById("b_cookie_consent")
                    .removeAttribute("disabled")
            } else {
                document.getElementById("b_cookie_consent")
                    .setAttribute("disabled", "true")
                }
    })
}

function setCookie() {
    if (getCookie("diceroller_user_id") == "") {
        document.cookie = "diceroller_user_id=0; path=/; secure=true"
    }   
    return true;
}

function checkCookie() {
    if (getCookie("diceroller_user_id") == "") {
        return false
    } else {
        return true
    }
}

function getCookie(cname) {
    let name = cname + "=";
    let decodedCookie = decodeURIComponent(document.cookie);
    let ca = decodedCookie.split(';');
    for(let i = 0; i <ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
  }
  
  function smallScreen() {
    if (window.innerWidth < 800) {
        document.getElementById("i_display_dice").remove()
        document.getElementById("l_display_dice").remove()
    }

    if (navigator.userAgent.toLowerCase().match(/mobile/i)) {
        document.getElementById("p_share").remove()
        document.getElementById("i_display_dice").remove()
        document.getElementById("l_display_dice").remove()
        document.getElementById("l_highlight_own_rolls").remove()
   
    }
  }
  
  function setDice() {
    var sel = document.querySelector("#s_d20");
    const d20 = sel.options[sel.selectedIndex].value; 
    sel = document.querySelector("#s_d12");  
    const d12 = sel.options[sel.selectedIndex].value; 
    sel = document.querySelector("#s_d10");  
    const d10 = sel.options[sel.selectedIndex].value; 
    sel = document.querySelector("#s_d8");
    const d8 = sel.options[sel.selectedIndex].value; 
    sel = document.querySelector("#s_d6");
    const d6 = sel.options[sel.selectedIndex].value; 
    sel = document.querySelector("#s_d4");
    const d4 = sel.options[sel.selectedIndex].value; 
    
    var first = true
    var isEmpty = false
    var text= "["
    for (let i = 0; i < d20; i++) {
        if (first) {
            first = false
        } else {
            text += ", "
        }
        text += "20"
        isEmpty = false
    }
    for (let i = 0; i < d12; i++) {
        if (first) {
            first = false
        } else {
            text += ", "
        }
        text += "12"
        isEmpty = false
    }
    for (let i = 0; i < d10; i++) {
        if (first) {
            first = false
        } else {
            text += ", "
        }
        text += "10"
        isEmpty = false
    }
    for (let i = 0; i < d8; i++) {
        if (first) {
            first = false
        } else {
            text += ", "
        }
        text += "8"
        isEmpty = false
    }
    for (let i = 0; i < d6; i++) {
        if (first) {
            first = false
        } else {
            text += ", "
        }
        text += "6"
        isEmpty = false
    }
    for (let i = 0; i < d4; i++) {
        if (first) {
            first = false
        } else {
            text += ", "
        }
        text += "4"
        isEmpty = false
    }
    text += "]"
    return text
}

function resetDice() {
    let elements =  [
        document.querySelector("#s_d20"),
        document.querySelector("#s_d12"),
        document.querySelector("#s_d10"),
        document.querySelector("#s_d8"),
        document.querySelector("#s_d6"),
        document.querySelector("#s_d4"),

        document.querySelector("#s_d20"),
        document.querySelector("#s_d12"),
        document.querySelector("#s_d10"),
        document.querySelector("#s_d8"),
        document.querySelector("#s_d6"),
        document.querySelector("#s_d4")
    ]
    let e = new Event("change");
    for (let i = 0; i < elements.length; i++) {
        elements[i].value = 0;
        elements[i].dispatchEvent(e);
    }
}


function hideAllDiceForm() {
    resetDice();                    
    document.querySelector(".all_dice_form")
        .classList.add("hidden")
}
 
  
  function initAllDiceForm() {
    let element = document.getElementById("f_roll_all");
    if (element == null) {
        return
    }

    element.addEventListener("submit", (event) => {
        event.preventDefault();
        const loc = location.href
        const player_id = "0"
        const dice = setDice(); 
        const chr = document.getElementById("f_name").value;
        const action = document.getElementById("f_action").value;

        fetch(loc, {
            method: "POST",
            headers: {
                "contentType": "application/json"
            },
            body: JSON.stringify({
                "dice": dice,
                "char": chr,
                "action": action
            })
        })
        .then( (response) => {
            if (response.ok) {
                hideAllDiceForm()
            }
        })
    });
       
    document.getElementById("b_close_all_dice")
        .addEventListener("click", (event) => {
            hideAllDiceForm()
        });
    document.getElementById("b_show_all_dice")
        .addEventListener("click", (event) => {
            document.querySelector(".all_dice_form")
                .classList.remove("hidden")
        });

    
  }