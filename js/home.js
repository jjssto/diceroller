const buttons = document.querySelectorAll(".b_room-selection");

buttons.forEach(element => {
    element.addEventListener("click", function(event) {
        var val = element.value;
        fetch("/", {
            method: "POST",
            headers: {
                "Accept": "application/json",
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                "id": val 
            }) 
        })
        .then((response) => {
            if (response.ok) {
                response.json().then(
                    data => {
                        const loc = `${window.location.href}room/${data}`;
                        window.location.assign(loc);
                    }
                );
            }
        })
    })
});