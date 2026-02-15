{% extends "blank.tpl" %}
{% block body %}
<!-- =========={ 404 Error Page }==========  -->
<div id="error-404" class="relative z-0 pt-20 lg:pt-24 pb-20 lg:pb-32 text-gray-300 bg-indigo-600 bg-gradient-to-r from-indigo-600 via-indigo-500 to-teal-500 dark:from-gray-800 dark:via-gray-700 dark:to-green-700 overflow-hidden min-h-screen flex items-center">
    <!-- Animated background elements -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
        <div class="absolute top-1/4 left-1/4 w-64 h-64 bg-white/5 rounded-full blur-3xl animate-pulse"></div>
        <div class="absolute bottom-1/4 right-1/4 w-96 h-96 bg-teal-300/5 rounded-full blur-3xl animate-pulse" style="animation-delay: 1s;"></div>
    </div>

    <div class="container xl:max-w-6xl mx-auto px-4 relative z-10">
        <!-- row -->
        <div class="flex flex-wrap flex-row -mx-4 justify-center">
            <!-- error content -->
            <div class="flex-shrink max-w-full px-4 w-full lg:w-2/3 self-center">
                <div class="text-center">
                    <!-- Decorative Icon -->
                    <div class="mb-8 inline-block">
                        <svg class="w-32 h-32 lg:w-40 lg:h-40 mx-auto mb-6 opacity-80 animate-bounce" style="animation-duration: 3s;" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                        </svg>
                    </div>

                    <!-- Main Error Content -->
                    <div class="mb-12">
                        <h1 class="text-8xl md:text-9xl lg:text-[12rem] leading-none mb-6 font-bold tracking-tight animate-pulse" style="text-shadow: 0 0 40px rgba(255,255,255,0.2);">404</h1>
                        <h2 class="text-3xl md:text-4xl lg:text-5xl leading-tight mb-6 font-bold text-white">Page Not Found</h2>
                        <div class="max-w-2xl mx-auto mb-8">
                            <p class="text-lg lg:text-xl leading-relaxed font-light text-gray-100 mb-6">
                                {% if msg %}
                                    {{ msg }}
                                {% else %}
                                    Oops! The page you're looking for doesn't exist. It might have been moved or deleted.
                                {% endif %}
                            </p>
                        </div>
                    </div>

                    <!-- Action Buttons -->
                    <div class="mb-12">
                        <div class="flex flex-wrap gap-4 justify-center mb-8">
                            <a class="group py-4 px-8 inline-flex items-center text-center rounded-lg leading-5 text-gray-900 bg-white border border-white hover:bg-gray-100 hover:border-gray-100 focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-indigo-600 transition-all duration-200 shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 font-semibold" href="/">
                                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" class="ltr:mr-2 rtl:ml-2 group-hover:scale-110 transition-transform" viewBox="0 0 16 16">
                                    <path d="M8.707 1.5a1 1 0 0 0-1.414 0L.646 8.146a.5.5 0 0 0 .708.708L2 8.207V13.5A1.5 1.5 0 0 0 3.5 15h9a1.5 1.5 0 0 0 1.5-1.5V8.207l.646.647a.5.5 0 0 0 .708-.708L13 5.793V2.5a.5.5 0 0 0-.5-.5h-1a.5.5 0 0 0-.5.5v1.293L8.707 1.5ZM13 7.207V13.5a.5.5 0 0 1-.5.5h-9a.5.5 0 0 1-.5-.5V7.207l5-5 5 5Z"/>
                                </svg>
                                Back to Home
                            </a>
                        </div>
                    </div>
                </div>
            </div>
        </div><!-- end row -->
    </div>
</div><!-- end 404 error -->

<style>
@keyframes float {
    0%, 100% { transform: translateY(0px); }
    50% { transform: translateY(-20px); }
}
.animate-bounce {
    animation: float 3s ease-in-out infinite;
}
</style>
{% endblock %}
