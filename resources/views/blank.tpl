<!DOCTYPE html>
<html lang="en" dir="ltr">
<head>
    <!-- Required meta tags -->
    <meta charset="UTF-8"/>
    <meta http-equiv="X-UA-Compatible" content="IE=edge"/>
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no"/>

    <!-- Title  -->
    <title>{{ title_page }}</title>
    <meta name="description" content="Built on top of FastHttp, the fastest HTTP engine for Go. Quick development with zero memory allocation and high performance. Very simple and easy to use."/>

    <!-- Development css (used in all pages) -->
    <link rel="stylesheet" id="stylesheet" href="/css/style.css"/>

    <!-- google font -->
    <link href="https://fonts.googleapis.com/css2?family=Nunito:wght@300;400;600;700&amp;display=swap" rel="stylesheet"/>

    <!-- Favicon  -->
    <link rel="icon" href="/assets/favicon.png"/>
</head>
<body class="font-sans text-base font-normal text-gray-600 dark:text-gray-400 dark:bg-gray-900 pt-18">

<!-- =========={ MAIN }==========  -->
<main id="content">
    {% block body %}{% endblock %}
</main><!-- end main -->

<script src="/js/alpine.min.js"></script><!-- core js -->
</body>
</html>
