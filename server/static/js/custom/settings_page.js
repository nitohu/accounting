function generateAPIKey() {
    let key = ""
    let seed = String(Date.now())
    // With the extra - it will be 64 characters
    let max = 61
    for(let i = 0; i < max; i++) {
        if (i == Math.floor(max*(1/4)) || i == Math.floor(max*(2/4)) || i == Math.floor(max*(3/4))) {
            key += "-"
        }
        if (i % 2 == 0 && (i/2) < seed.length) {
            key += seed[i/2]
        } else {
            let n = Math.random() * 80 + 46
            key += String.fromCharCode(n)
        }
    }
    document.getElementById("api_key").value = key
}