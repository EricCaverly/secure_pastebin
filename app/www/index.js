
// Effectively ripped from https://dev.to/shubhamkhan/beginners-guide-to-aes-encryption-and-decryption-in-javascript-using-cryptojs-592
function enc_message(message, psk) {
    let key = CryptoJS.SHA256(psk).toString();

    // Generate a random Initialization Vector (IV) for security
    const iv = CryptoJS.lib.WordArray.random(16);

    // Encrypt the text using AES with CBC mode and the secret key
    const encrypted = CryptoJS.AES.encrypt(message, CryptoJS.enc.Hex.parse(key), {
        iv: iv,
        padding: CryptoJS.pad.Pkcs7,
        mode: CryptoJS.mode.CBC,
    });

    // Concatenate IV and ciphertext and encode in Base64 format
    const encryptedBase64 = CryptoJS.enc.Base64.stringify(
        iv.concat(encrypted.ciphertext)
    );

    return encryptedBase64;
}


// Effectively ripped from https://dev.to/shubhamkhan/beginners-guide-to-aes-encryption-and-decryption-in-javascript-using-cryptojs-592
function dec_message(encText, psk) {
    let key = CryptoJS.SHA256(psk).toString();

    try {
        const fullCipher = CryptoJS.enc.Base64.parse(encText);

        // Extract IV and ciphertext from the parsed cipher
        const iv = CryptoJS.lib.WordArray.create(fullCipher.words.slice(0, 4), 16);
        const ciphertext = CryptoJS.lib.WordArray.create(fullCipher.words.slice(4));

        const cipherParams = CryptoJS.lib.CipherParams.create({
            ciphertext: ciphertext,
        });

        // Decrypt the ciphertext using AES and the provided secret key
        const decrypted = CryptoJS.AES.decrypt(cipherParams, CryptoJS.enc.Hex.parse(key), {
            iv: iv,
            padding: CryptoJS.pad.Pkcs7,
            mode: CryptoJS.mode.CBC,
        });

        // Return decrypted text in UTF-8 format
        return decrypted.toString(CryptoJS.enc.Utf8);
        
    } catch (error) {
        console.error("decryption error:", error);
        return null;
    }
}


function api_req(method, endpoint, data, success_func) {
    const opt = {
        url: `/api/${endpoint}`,
        type: method,
        data: data
    }

    let req_obj = $.ajax(opt);

    req_obj.fail((xhr_err, _, err) => {
        console.log(xhr_err);
        console.log(err);
        alert(`There was a problem making the request`);
    });

    req_obj.done(success_func);
}


function setup_note_creation() {
    const form = $("#note_form");
    const loading = $("#loading_card");
    const create_note = $("#create_card");
    const result_card = $("#result_card");
    const result_body = $("#result_body");

    // No need to fetch any data right away since we are making the note
    form.submit((e) => {
        e.preventDefault();

        let msg = $("#new_content").val();
        let psk = $("#new_passphrase").val();
        let ipr = $("#new_ip_restriction").val();

        let ciphertext = enc_message(msg, psk);

        create_note.hide();
        loading.show();

        api_req("POST", "note", {
            "content": ciphertext, 
            "allowed_ips": ipr,
        }, (result) => {
            console.log(result);
            loading.hide();
            result_card.show();
            
            if (result.success) {
                const url = `${window.location.href}?uuid=${result.data}`;

                const btn = document.createElement("button");
                btn.setAttribute("class", "btn btn-primary");
                btn.innerHTML = `&#x1F4CB;`;
                btn.addEventListener("click", () => {
                    navigator.clipboard.writeText(url)
                });

                result_body.empty();
                result_body.append(`&#x2713; Note available <a href="${url}">here</a> `);
                result_body.append(btn);
            } else {
                result_body.html(`&#x274c; Error: ${result.message}`);
            }
        });
    })

    loading.hide();
    create_note.show();
}


function setup_note_retrieval(uuid) {
    const loading = $("#loading_card");
    const dec_card = $("#decrypt_card");
    const result_card = $("#result_card");
    const result_body = $("#result_body");

    // Obtain the note, rendering a field for the passphrase if the note eixsts, or an error
    api_req("GET", `note/${uuid}`, {}, (result) => {
        if (result.success) {
            $("#decrypt_form").submit((e) => {
                e.preventDefault();
                loading.show();
                
                let psk = $("#view_passphrase").val();

                let msg = dec_message(result.data.content, psk);
                console.log(msg);

                result_body.empty();
                result_body.html(document.createTextNode(msg));

                loading.hide();
                result_card.show();
            });
            loading.hide();
            dec_card.show();
        } else {
            loading.hide();
            result_card.show();
            result_body.html(`&#x274c; Error: ${result.message}`);
        }
    });
}


$(() => {
    // Check if the UUID is specified as a Query Parameter
    const params = new Proxy(new URLSearchParams(window.location.search), {
        get: (searchParams, prop) => searchParams.get(prop),
    });
    let uuid = params.uuid;

    // Render UI accordingly
    if (uuid == null) {
        setup_note_creation();
    } else {
        setup_note_retrieval(uuid);
    }  
});