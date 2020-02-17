<!DOCTYPE html>
<html>

<body>
  <head>
    <link rel='stylesheet' href='styles.css'>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css">
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.4.1/jquery.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.16.0/umd/popper.min.js"></script>
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.4.1/js/bootstrap.min.js"></script>
    <meta name="viewport" content="width=device-width, initial-scale=1" charset='utf-8'>
  </head>
  <nav class="navbar navbar-expand-sm navbar-custom">
    <img src='/Logos/Logo2.gif' id='logo' onclick='document.location.href="/"'>
    <ul class="navbar-nav">
      <li class="nav-item">
        <a class="nav-link" href="/">Home</a>
      </li>
      <li class="nav-item">
        <a class="nav-link" href="/docs/intro.html">Documentation</a>
      </li>
      <li class="nav-item">
        <a class="nav-link" href="/downloads.php">Downloads</a>
      </li>
      <li class="nav-item">
        <a class="nav-link" href="/about.php">About</a>
      </li>
    </ul>
  </nav>
  <p id='about-html'>
    Omm is a general purpose mathematics programming language created by me, Ankit Karnam, and concieved in October of 2019. I created Omm because I came up with an idea for near infinite mathematics and precision. The syntax I chose was very different from ordinary, because I remove all whitespace before parsing the code. When I came up with the syntax, I wanted to maintain a similar syntax to javascript, but simplify it a bit. Omm is an object-oriented and logic-oriented programming language. Omm is a very versitale language, as it runs javascript code inside of it. I will be working to add support for more languages inside Omm.
  </p>
  <div class='footer'>
    <img src='/Logos/Logo2.gif' id='foot-img' onclick='document.location.href="/"'>
    <a id='foot-img-txt'>Omm</a>
    <br>
    <a class='footer-copyright textcenter py-3' id='copyright'>© 2019 Ankit Karnam</a>
    <br>
    <small id='license'>Omm is under the <a href='/license.txt'>MIT</a> license</small>
    <div class='foot-ext-div'>
      <a href='https://github.com/Ankizle/Omm/issues' class='foot-ext'>Issues</a>
    </div>
  </div>
</body>

</html>
