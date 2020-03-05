$(document).ready(function() {
    let btns = document.getElementsByClassName("deleteEntry")

    for (let i = 0; i < btns.length; i++) {
        let btn = btns[i]
        btn.addEventListener("click", deleteAccount)
    }

    function deleteAccount(e) {
        data = {
            "ID": Number.parseInt(e.path[2].id)
        }
        console.log(e.path)
        console.log(e.path[2].id)
        console.log(typeof(e.path[2].id))
        console.log(data)
        let xhr = new XMLHttpRequest()
        xhr.open("DELETE", "/api/accounts/delete", true)
        xhr.setRequestHeader("Content-Type", "application/json")
        xhr.onreadystatechange = function() {
            if (this.readyState == 4) {
                if (this.status == 400) {
                    console.error(this.response)
                } else if (this.status == 200) {
                    window.location.reload()
                } else {
                    console.warn(this.status)
                }
            }
        }
        xhr.send(JSON.stringify(data))
    }
})
