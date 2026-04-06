function getData() {
  // get value from id in html
  let name = document.getElementById("name").value;
  let email = document.getElementById("email").value;
  let phoneNumber = document.getElementById("phone-number").value;
  let subject = document.getElementById("subject").value;
  let msg = document.getElementById("msg").value;

  if ((name, email, phoneNumber, subject, (message = ""))) {
    alert("Tolong data segera diisi");
  } else {
    alert("Data akan segera dikirim melalui email");
  }

  const emailReceive = "hanseenn01@gmail.com";

  let mailTo = document.createElement("a");
  mailTo.href = `mailto:${emailReceive}?subject=${subject}&body=Halo nama saya ${name}, saya ingin ${msg}, bisakah anda menghubungi saya di ${phoneNumber}`;
  mailTo.click();

  let show = {
    name,
    email,
    phoneNumer,
    subject,
    msg,
  };
  console.log(show);
}