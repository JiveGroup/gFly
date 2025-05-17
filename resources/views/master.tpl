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
    <!-- =========={ Header Menu }==========  -->
    {% if account != nil %}
    <header class="fixed top-0 left-0 right-0 z-50 bg-indigo-700 text-white py-4">
        <div class="container xl:max-w-6xl mx-auto px-4">
            <div class="flex justify-between items-center">
                <div>
                    <a href="/" class="text-white font-bold hover:text-gray-200 mr-6">Home</a>
                </div>
                <div>
                    <a href="/profile" class="text-white font-bold hover:text-gray-200 mr-6">Profile</a>
                </div>
                {% if !isPaths("/") %}
                <div>
                    <a href="/users" class="text-white font-bold hover:text-gray-200 mr-6">Users</a>
                </div>
                <div>
                    <button id="logout-btn" class="text-white font-bold hover:text-gray-200">Logout</button>
                </div>
                {% endif %}
            </div>
        </div>
    </header>
    {% endif %}
    {% if account == nil %}
    <header class="fixed top-0 left-0 right-0 z-50 bg-indigo-700 text-white py-4">
        <div class="container xl:max-w-6xl mx-auto px-4">
            <div class="flex justify-between items-center">
                <div>
                    <a href="/" class="text-white font-bold hover:text-gray-200 mr-6">Home</a>
                </div>
                {% if !isPaths("/login") %}
                <div>
                    <a href="/login" class="text-white font-bold hover:text-gray-200 mr-6">Login</a>
                </div>
                {% endif %}
            </div>
        </div>
    </header>
    {% endif %}
    {% block body %}{% endblock %}
</main><!-- end main -->

<script src="/vendors/alpinejs/dist/cdn.min.js"></script><!-- core js -->

</body>
</html>
