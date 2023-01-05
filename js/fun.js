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
    initReset,
    createRowDice,
    setResultHeaderDice
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
    var offsetStr = sessionStorage.getItem("ts_offset");
    if (offsetStr == null || offsetStr.length == 0) {
        const date = new Date();
        const offset = -1 * date.getTimezoneOffset() * 60; 
        offsetStr = offset.toString();
        sessionStorage.setItem("ts_offset", offsetStr)
    }
    if (last_roll != "") {
        target += "/" + last_roll
    }
    fetch(target, {
        method: "GET",
        headers: {
            "ts_offset": offsetStr
        },
    })
    .then(response => {
        response.text().then(data => {
                insertRolls(createRow, data);
        })
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

function init(createRow) {

    window.addEventListener("DOMContentLoaded", () => {
        setHighlightOwnRolls()
        setDisplayDice()
        setColor()
        detectFocus()
        settingVisibility()
        roomSettingForm()
        setLink()

        if (displayDice) {
            window.setInterval(() => {getRolls(createRowDice)}, 1000);
        } else {
            window.setInterval(() => {getRolls(createRow)}, 1000);
        }
    })
    window.addEventListener("load", () => {
        sessionStorage.setItem("last_roll", "");
        if (displayDice) {
            getRolls(createRowDice);
        } else {
            getRolls(createRow);
        }
    });
}

function initReset(reset) {
    document.getElementById("b_reset").addEventListener("click", () => {
        reset();
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
    