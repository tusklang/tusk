<!DOCTYPE html>
<html>

<body>
  <head>

    <meta charset='utf-8'>

    <!-- Semantic-UI stuff -->
    <link rel='stylesheet' href='/semantic/semantic.min.css'>
    <script src='https://code.jquery.com/jquery-3.1.1.min.js'></script>
    <script src='semantic/semantic.min.js'></script>
    <!----------------------->

    <!-- Font Awesome stuff -->
    <script src='https://kit.fontawesome.com/3df27ff5ea.js' crossorigin='anonymous'></script>
    <!------------------------>

    <!-- Favicon -->
    <link rel='shortcut icon' href='favicon.ico'>
    <!------------->

    <script src='include.js'></script>
    <link rel='stylesheet' href='styles.css'>
  </head>
  
  <div id='prenav'></div>

  <div class='ui grid' id='index-body'>
    <div class='ui container'>
      <h1 id='omm'>
        Omm
      </h1>

      <p id='omm-body'>
        Omm is a general purpose language with arbitrary precision
      </p>

      <button class='ui download-btn button' onclick='download()'>Download Latest Version (Beta 1.0.0)</button>

      <script>

        function download() {

          switch (navigator.platform) {
            case "Win32":
              window.location.href = '/versions/beta1.0.0/setup.msi'
              break;
            default:
              alert("Sorry, but the Omm installer is not available on your platform. Try installing from the source")
          }

        }

      </script>

      <div id='index-background-div' class='ui large disabled medium right floated image'>
        <img id='index-background' src='Logos/in-use/index-background.png'>
      </div>
    </div>
  </div>

  <div id='features'>

    <div class='ui segment inverted feature'>

      <h1 class='feature-head'>Object Oriented</h1>
      <div class='ui horizontal divider'>
        <i class='fas fa-copy fa-2x feature-ico'></i>
      </div>

      <p class='description'>
        Omm is an object oriented language, and can create a more scalable and re-usable infrastructure.
      </p>

    </div>

    <div class='ui segment inverted feature'>

      <h1 class='feature-head'>Precision</h1>
      <div class='ui horizontal divider'>
        <i class='fas fa-calculator fa-2x feature-ico'></i>
      </div>

      <p class='description'>
        Omm has arbitrary integer and decimal precision. This means that Omm is not capped by the 64 bit limit, but rather the system's memory.
      </p>

    </div>

    <div class='ui segment inverted feature'>

      <h1 class='feature-head'>Security</h1>
      <div class='ui horizontal divider'>
        <i class='fas fa-lock fa-2x feature-ico'></i>
      </div>

      <p class='description'>
        Omm has private and public variables, so you don't need to mess with currying to secure code.
      </p>

    </div>

    <div class='ui segment inverted feature'>

      <h1 class='feature-head'>Learning</h1>
      <div class='ui horizontal divider'>
        <i class='fas fa-brain fa-2x feature-ico'></i>
      </div>

      <p class='description'>
        Omm has easy and intuitive syntax which makes it easy to learn.
      </p>

    </div>

  </div>

  <div id='prefoot'></div>

</body>

</html>
