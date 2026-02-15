{% extends "master.tpl" %}
    {% block body %}
<!-- =========={ Hero }==========  -->
<div id="hero" class="relative z-0 pt-36 lg:pt-44 xl:pt-48 pb-20 lg:pb-32 text-gray-300 bg-indigo-600 bg-gradient-to-r from-indigo-600 via-indigo-500 to-teal-500 dark:from-gray-800 dark:via-gray-700 dark:to-green-700 overflow-hidden h-screen">
    <div class="container xl:max-w-6xl mx-auto px-4">
        <!-- row -->
        <div class="flex flex-wrap flex-row -mx-4 justify-center">
            <!-- hero content -->
            <div class="flex-shrink max-w-full px-4 w-full md:w-9/12 lg:w-1/2 self-center lg:ltr:pr-12 lg:rtl:pl-12">
                <div class="text-center lg:ltr:text-left lg:rtl:text-right mt-6 lg:mt-0">
                    <div class="mb-8">
                        <h1 class="text-4xl lg:text-5xl leading-normal mb-3 font-bold">{{ hero_text }}</h1>
                        <p class="600 leading-relaxed font-light text-xl mx-auto pb-2">
                            Built on top of <a href="https://github.com/valyala/fasthttp" target="_blank" class="font-bold" style="color: var(--white)">FastHttp - the fastest HTTP engine</a> and <a href="https://github.com/jivegroup/fluentsql" target="_blank" class="font-bold" style="color: var(--white)">FluentSQL - flexible and powerful SQL builder</a> for Go. Quick development with minimum memory allocation and high performance. Very simple and easy to use.
                        </p>
                    </div>
                    <a class="py-2 px-4 inline-block text-center rounded leading-5 text-gray-700 bg-gray-300 border border-gray-300  hover:bg-gray-200 hover:ring-0 hover:border-gray-200 focus:bg-gray-200 focus:border-gray-200 focus:outline-none focus:ring-0 mr-4" target="_blank" href="https://gfly.dev">
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="inline-block ltr:mr-2 rtl:ml-2 bi bi-check2-circle" viewBox="0 0 16 16">
                            <path d="M2.5 8a5.5 5.5 0 0 1 8.25-4.764.5.5 0 0 0 .5-.866A6.5 6.5 0 1 0 14.5 8a.5.5 0 0 0-1 0 5.5 5.5 0 1 1-11 0z"/>
                            <path d="M15.354 3.354a.5.5 0 0 0-.708-.708L8 9.293 5.354 6.646a.5.5 0 1 0-.708.708l3 3a.5.5 0 0 0 .708 0l7-7z"/>
                        </svg>Visit Us
                    </a>
                </div>
            </div>

            <!-- hero image -->
            <div class="flex-shrink max-w-full px-4 w-full md:w-9/12 lg:w-1/2 self-center">
                <div class="px-12 md:ml-16 md:pr-0 mt-4">
                    <img src="assets/hero.png" alt="Hero Image" class="max-w-full mx-auto"/>
                </div>
            </div>
        </div><!-- end row -->
    </div>
</div><!-- end hero -->
    {% endblock %}
