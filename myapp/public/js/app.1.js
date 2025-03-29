function resetSearch() {
    let el = document.getElementById('filter');
    el.value='';
    el.focus();

    document.querySelector('.grid-scroll').scrollTop=0;
}

// Example starter JavaScript for disabling form submissions if there are invalid fields
function validate () {
    // Fetch all the forms we want to apply custom Bootstrap validation styles to
    const forms = document.querySelectorAll('.needs-validation');

    // Loop over them and prevent submission
    forms.forEach(function (form) {
        form.addEventListener('submit', function (event) {
            if (!form.checkValidity()) {
                event.preventDefault();
                event.stopPropagation();
            }

            form.classList.add('was-validated');
        }, false);
    });

    document.body.addEventListener('htmx:configRequest', function (event) {
        const form = event.target.closest('form');
        if (form && !form.checkValidity()) {
            event.preventDefault();
            form.classList.add('was-validated');
        }
    });
}

function notify(msg, msgType) {
    Swal.fire({
        position: "top-end",
        icon: msgType,
        title: msg,
        showConfirmButton: false,
        timer: 3000,
        toast: true
    });

    console.log("msg", msg)
}
