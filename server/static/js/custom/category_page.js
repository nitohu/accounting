$("document").ready(function() {
    var categoryID = 0
    var tbody = document.getElementById("categoryTable")

    // debounce(getCategories, 200, false)()
    getCategories()
    let addLineBtn = document.getElementById("addLineBtn")
    addLineBtn.addEventListener("click", function() {
        categoryID = 0
        console.log("Add a line")
    })

    let formBtn = document.getElementById("formBtn")
    formBtn.addEventListener("click", saveCategories)

    function debounce(func, wait, immediate) {
        var timeout;

        return function executedFunction() {
            var context = this;
            var args = arguments;
                
            var later = function() {
            timeout = null;
            if (!immediate) func.apply(context, args);
            };

            var callNow = immediate && !timeout;
            
            clearTimeout(timeout);

            timeout = setTimeout(later, wait);
            
            if (callNow) func.apply(context, args);
        };
    };

    function saveCategories() {
        console.log("saveCategories()")

        let name = document.getElementById("name").value
        let hex  = document.getElementById("hex").value

        let data = {
            "ID": categoryID,
            "Name": String(name),
            "Hex": String(hex),
        }

        let xhr = new XMLHttpRequest()
        xhr.open("POST", "/api/categories/create")
        xhr.setRequestHeader("Content-Type", "application/json")
        xhr.onreadystatechange = debounce(getCategories, 200, false)
        xhr.onload = () => {
            console.log(xhr.status)
            console.log(xhr.response)
        }
        console.log(data)
        xhr.send(JSON.stringify(data))
    }

    function getCategories() {
        console.log("getCategories()")
    
        let data = {
            "ID": 0,
        }
        let xhr = new XMLHttpRequest()
        xhr.open("GET", "/api/categories", true)
        xhr.setRequestHeader("Content-Type", "application/json")
        xhr.onreadystatechange = function() {
            // Remove current categories in DOM
            while (tbody.lastElementChild) {
                tbody.removeChild(tbody.lastElementChild)
            }
            // Append new categories
            if (this.readyState == 4 && this.status == 200) {
                let categories = JSON.parse(this.responseText)
                for(let i in categories) {
                    let category = categories[i]
                    console.log(category["Name"])
                    let content = document.createElement("tr")
                    content.setAttribute("style", 'background-color: '+category["Hex"])
                    content.setAttribute("class", "categoryItem")

                    let td = document.createElement("td")
                    td.innerHTML = category["Hex"]
                    content.appendChild(td)

                    td = document.createElement("td")
                    let a = document.createElement("a")
                    a.setAttribute("href", "#")
                    a.setAttribute("data-toggle", "modal")
                    a.setAttribute("data-target", "#formModal")
                    a.innerHTML=category["Name"]
                    td.appendChild(a)
                    content.appendChild(td)

                    td = document.createElement("td")
                    td.innerHTML=category["CreateDate"]
                    content.appendChild(td)

                    td = document.createElement("td")
                    td.innerHTML=category["LastUpdate"]
                    content.appendChild(td)

                    td = document.createElement("td")
                    td.innerHTML=category["TransactionCount"]
                    content.appendChild(td)

                    td = document.createElement("td")
                    td.setAttribute("class", "deleteEntry")
                    a = document.createElement("i")
                    a.setAttribute("class", "zmdi zmdi-close")
                    td.appendChild(a)
                    content.appendChild(td)
                    
                    // let row = document.getElementById("addLineRow")
                    // tbody.insertBefore(content, row)
                    tbody.appendChild(content)
                }
            } else {
                console.warn(this.status)
                console.warn(this.readyState)
                console.log(this.response)
            }
        }
        xhr.send(JSON.stringify(data))
    }

    // let btns = document.getElementsByClassName("deleteEntry")

    // for (let i = 0; i < btns.length; i++) {
    //     let btn = btns[i]
    //     btn.addEventListener("click", function() {
    //         let id = btn.id.split("_")[1]

    //         console.warn("Deleting record with id ", id)
    //         fetch("/categories/delete/"+id).then(function() {
    //             window.location.reload()
    //         })

    //     })
    // }
})