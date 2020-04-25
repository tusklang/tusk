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

    <script src='curtab.js'></script>
    <link rel='stylesheet' href='styles.css'>
  </head>

  <div id='navbar' class='ui inverted labeled icon large eight item menu'>

    <a class='item menu-item' href='/index.php'>
      <i class='fas fa-home fa-2x'></i>
      Home
    </a>
    <a class='item menu-item' href='/downloads.php'>
      <i class='fas fa-download fa-2x'></i>
      Downloads
    </a>
    <a class='item menu-item' href='/about.php'>
      <i class='fas fa-info fa-2x'></i>
      About
    </a>
    <a class='item menu-item' href='/contributing.php'>
      <i class='fas fa-hands-helping fa-2x'></i>
      Contributing
    </a>
    <a class='item menu-item' href='/donations.php'>
      <i class='fas fa-donate fa-2x'></i>
      Donations
    </a>
    <a class='item menu-item' href='/license.php'>
      <i class='fas fa-balance-scale fa-2x'></i>
      License
    </a>
    <a class='item menu-item' href='/docs.php'>
      <i class='fas fa-book-open fa-2x'></i>
      Documentation
    </a>
    <a class='item menu-item' href='https://github.com/Ankizle/Omm'>
      <i class='fab fa-github fa-2x'></i>
      Github
    </a>

  </div>

  <div class='ui grid' id='index-body'>
    <div class='ui container'>
      <h1 id='omm'>
        Omm
      </h1>

      <p id='omm-body'>
        Omm is a general purpose language with arbitrary precision
      </p>

      <button class='ui download-btn button'>Download Latest Version (v1.0.0)</button>

      <div id='index-background-div' class='ui large disabled medium right floated image'>
        <img id='index-background' src='Logos/in-use/index-background.png'>
      </div>
    </div>
  </div>

  <div id='feature' class='ui one center aligned page grid'>
    <div class='ui segment inverted'>

      <h1 class='feature-head'>Compatiblility</h1>
      <div class='ui horizontal divider'>
        <i class='fas fa-puzzle-piece fa-2x feature-ico'></i>
      </div>

      <p class='description'>
        Omm has an extension named OES that allows you to embed Omm scripts within other languages.
      </p>

    </div>

    <div class='ui segment inverted'>

      <h1 class='feature-head'>Precision</h1>
      <div class='ui horizontal divider'>
        <i class='fas fa-calculator fa-2x feature-ico'></i>
      </div>

      <p class='description'>
        Omm, like python, has arbitrary precision. This means that the maximum number is limited by the machine's memory, rather than the 64 bit maxint
      </p>

    </div>

    <div class='ui segment inverted'>

      <h1 class='feature-head'>Ease Of Learning</h1>
      <div class='ui horizontal divider'>
        <i class='fas fa-brain fa-2x feature-ico'></i>
      </div>

      <p class='description'>
        Omm is a very easy language to understand, so it can be easily integrated into your project.
      </p>

    </div>
  </div>
</body>

</html>
