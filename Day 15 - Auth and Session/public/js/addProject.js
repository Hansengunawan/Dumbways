let blogs = [];

function getProject(e) {
  e.preventDefault();

  //get id from html
  let projectName = document.getElementById("projectName").value;
  let startDate = document.getElementById("startDate").value;
  let endDate = document.getElementById("endDate").value;
  let desc = document.getElementById("desc").value;
  let upImg = document.getElementById("upFile").files;

  // declare icon technologies
  const tech1 = '<i class="fa-brands fa-node"></i>';
  const tech2 = '<i class="fa-brands fa-react"></i>';
  const tech3 = '<i class="fa-brands fa-python"></i>';
  const tech4 = '<i class="fa-brands fa-golang"></i>';

  //check icon
  let nodeJS = document.getElementById("tech-one").checked ? tech1 : "";
  let reactJS = document.getElementById("tech-two").checked ? tech2 : "";
  let python = document.getElementById("tech-three").checked ? tech3 : "";
  let golang = document.getElementById("tech-four").checked ? tech4 : "";

  //validation
  if (projectName == "") {
    return alert("Tolong isi bagian nama projek anda");
  } else if (startDate && endDate == "") {
    return alert("Tolong masukkan tanggal yang sesuai");
  } else if (desc == "") {
    return alert("Tolong masukkan deskripsi dengan benar");
  } else if (nodeJS && reactJS && python && golang == "") {
    return alert("Tolong pilh salah satu");
  } else {
    alert("Segera ditampilkan");
  }

  // convert image to blob object
  upImg = URL.createObjectURL(upImg[0]);

  startDate = new Date(startDate);
  endDate = new Date(endDate);

  const distance = endDate - startDate;

  //count duration
  let duration = Math.floor(distance / (12 * 30 * 7 * 24 * 60 * 60 * 1000));
  if (duration > 0) duration = `${duration} years`; //years
  else {
    duration = Math.floor(distance / (30 * 24 * 60 * 60 * 1000));
    if (duration > 0) duration = `${duration} months`; //month
    else {
      duration = Math.floor(distance / (7 * 24 * 60 * 60 * 1000));
      if (duration > 0) duration = `${duration} weeks`;
      else {
        duration = Math.floor(distance / (24 * 60 * 60 * 1000));
        if (duration > 0) duration = `${duration} days`;
        else {
          duration = Math.floor(distance / (60 * 60 * 1000));
          if (duration > 0) duration = `${duration} hours`;
          else {
            duration = Math.floor(distance / (60 * 1000));
            if (duration > 0) duration = `${duration} minutes`;
            else {
              duration = Math.floor(distance / 1000);
              if (duration > 0) duration = `${duration} seconds`;
            }
          }
        }
      }
    }
  }

  let dataProject = {
    projectName,
    startDate,
    endDate,
    duration,
    desc,
    upImg,
    nodeJS,
    reactJS,
    python,
    golang,
    author: "Hansen",
    postedAt: new Date(),
  };

  blogs.push(dataProject);
  console.log(blogs);
  clear();
  showProject();
}

function showProject() {
  document.getElementById("take").innerHTML = "";
  for (let i = 0; i <= blogs.length; i++) {
    document.getElementById(
      "take"
    ).innerHTML += `<div class="container-fluid px-5 my-lg-5 p-lg-5">
    <div class="row row-cols-1 row-cols-md-3 g-lg-5 mt-5">
        <div class="col">
            <div class="card h-100 p-lg-3" style="max-width:330px; box-shadow: 0 0 10px #dbdbdb;">
                <img src=${blogs[i].upImg} class="card-img" alt="image">
                <div class="card-body p-0 mt-4">
                    <a href="#" class="nav-link"><h6 class="card-title fw-bold">${blogs[i].projectName}</h6></a>
                    <small class="text-muted">${blogs[i].duration}</small>
                    <p class="card-text mt-3" style="text-align: justify;">${blogs[i].desc}</p>
                    <div class="card-text d-flex ">
                        <div class="mt-2 fs-2 pb-2 ">
                            <a class="navbar-brand me-4" href="#">
                            ${blogs[i].nodeJS}
                            </a>
                            <a class="navbar-brand me-4" href="#">
                            ${blogs[i].reactJS}
                            </a>
                            <a class="navbar-brand me-4 " href="#">
                            ${blogs[i].python}
                            </a>
                            <a class="navbar-brand me-4 " href="#">
                            ${blogs[i].golang}
                            </a>
                        </div>
                    </div>
                    <div class="card-text mt-2 d-flex">
                        <a class="btn bg-black text-white fw-semibold w-50 me-1 p-0 rounded-pill" href="#"
                            role="button">edit</a>
                        <a class="btn bg-black text-white w-50 ms-1 p-0 fw-semibold rounded-pill" href="#"
                            role="button">delete</a>
                    </div>
                </div>
            </div>
        </div>
`;
  }
}

function clear() {
  document.getElementById("projectName").value = "";
  document.getElementById("startDate").value = "";
  document.getElementById("endDate").value = "";
  document.getElementById("desc").value = "";
  document.getElementById("upFile").value = "";
}
