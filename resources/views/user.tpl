{% extends "master.tpl" %}
    {% block body %}
    <!-- =========={ User }==========  -->
<div id="hero" class="relative z-0 pt-36 lg:pt-44 xl:pt-48 pb-20 lg:pb-32 text-gray-300 bg-indigo-600 bg-gradient-to-r from-indigo-600 via-indigo-500 to-teal-500 dark:from-gray-800 dark:via-gray-700 dark:to-green-700 overflow-hidden h-screen">
    <div class="container xl:max-w-6xl mx-auto px-4">
        <div class="flex flex-wrap flex-row -mx-4 justify-center">
            <!-- Profile Box -->
            <div class="w-full max-w-md">
                <div class="dark:bg-gray-800 shadow-md rounded-lg px-8 py-6">
                    <h2 class="text-2xl font-bold mb-6 text-center">User list</h2>
                    {% include "user/list.tpl" %}
                </div>
            </div>
        </div>
    </div>
</div><!-- end user -->
    {% endblock %}
