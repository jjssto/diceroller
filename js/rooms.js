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


var buttons = document.querySelectorAll(".room_id_button")

for (var i = 0; i < buttons.length; i++) {
    buttons[i].addEventListener("click", (event) => {
        const url = "/room/" + event.target.value;
        window.open(url, '_blank');
    })
}

buttons = document.querySelectorAll(".delete_room_button")

for (var i = 0; i < buttons.length; i++) {
    buttons[i].addEventListener("click", (event) => {
        fetch("/rooms", {
            method: "POST",
            headers: {
                "contentType": "application/json",
                "accept": "application/json",
            },
            body: JSON.stringify({
                "room_id": event.target.value,
            })
        })
            .then(resp => {
                if (resp.ok) {
                    const name = "row_" + event.target.value;
                    var rows = document.getElementsByName(name);
                    for (var i = rows.length - 1; i >= 0; i--) {
                        rows[i].remove()
                    }
                } else {
                    window.alert("Deleting romm #" + event.target.value +
                        " failed!");
                }
            })
    })
}



document.getElementById("b_new").addEventListener("click", () => {
    window.open("/", "_self")
})

const date_cols = document.querySelectorAll(".date_col")
const comparison = '2000-01-01';

for (var i = 0; i < date_cols.length; i++) {
    const date = date_cols[i].textContent;
    if (date < comparison) {
        date_cols[i].textContent = ""
    }
}