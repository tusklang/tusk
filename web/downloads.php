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

  <table id='downloads' class='ui inverted black striped celled padded table'>
    <thead>
      <tr>
        <th>Version</th>
        <th>Description</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>
          <a id='download_version' onclick='downloadv("beta1.0.0")'>Beta 1.0.0</a>
        </td>
        <td>
          Omm Beta 1
        </td>
      </tr>
    </tbody>
  </table>

  <script>
    function downloadv(version) {
      switch (navigator.platform) {
        case "Win32":
          window.location.href = '/versions/' + version + '/setup.msi'
          break;
        default:
        alert("Sorry, but the Omm installer is not available on your platform. Try installing from the source")
      }
    }
  </script>

  <div id='prefoot'></div>

</body>

</html>
