{% extends "master.tpl" %}
    {% block body %}
    <!-- =========={ User }==========  -->
<div id="hero" class="relative z-0 pt-36 lg:pt-44 xl:pt-48 pb-20 lg:pb-32 text-gray-300 bg-indigo-600 bg-gradient-to-r from-indigo-600 via-indigo-500 to-teal-500 dark:from-gray-800 dark:via-gray-700 dark:to-green-700 overflow-hidden h-screen">
    <div class="container xl:max-w-6xl mx-auto px-4">
        <div class="flex flex-wrap flex-row -mx-4 justify-center">
            <!-- Login Form Card -->
            <div class="w-full max-w-md">
                <div class="bg-white dark:bg-gray-800 shadow-md rounded-lg px-8 py-6">
                    <h2 class="text-2xl font-bold text-gray-800 dark:text-white mb-6 text-center">Login</h2>
                    <form id="login-form" action="/login" method="post">
                        <!-- Email Input -->
                        <div class="mb-4">
                            <label class="block text-gray-700 dark:text-gray-300 text-sm font-bold mb-2" for="username">
                                Email
                            </label>
                            <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                   id="username"
                                   name="username"
                                   type="email"
                                   placeholder="your@email.com"
                                   value="admin@gfly.dev"
                                   required="required" />
                        </div>

                        <!-- Password Input -->
                        <div class="mb-6">
                            <label class="block text-gray-700 dark:text-gray-300 text-sm font-bold mb-2" for="password">
                                Password
                            </label>
                            <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                   id="password"
                                   name="password"
                                   type="password"
                                   value="P@seWor9"
                                   placeholder="••••••••"
                                   required="required" />
                        </div>

                        <!-- Remember Me Checkbox -->
                        <div class="mb-6 flex items-center">
                            <input id="remember" name="remember" type="checkbox" class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded">
                            <label for="remember" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
                                Remember me
                            </label>
                            <a href="/forgot-password" class="text-sm text-indigo-500 hover:text-indigo-700 ml-auto">
                                Forgot password?
                            </a>
                        </div>

                        <!-- Submit Button -->
                        <div class="flex items-center justify-between">
                            <button class="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline transition duration-150 ease-in-out" type="submit">
                                Sign In
                            </button>
                        </div>

                        <!-- Register Link -->
                        <div class="text-center mt-6">
                            <p class="text-sm text-gray-600 dark:text-gray-400">
                                Don't have an account?
                                <a href="/register" class="font-medium text-indigo-500 hover:text-indigo-700">
                                    Register
                                </a>
                            </p>
                        </div>
                    </form>

                    <script>
                        document.addEventListener('DOMContentLoaded', function() {
                            const form = document.getElementById('login-form');

                            form.addEventListener('submit', function(e) {
                                e.preventDefault();

                                const username = document.getElementById('username').value;
                                const password = document.getElementById('password').value;

                                // Make AJAX request to API
                                fetch('/api/v1/frontend/auth/signin', {
                                    method: 'POST',
                                    headers: {
                                        'Content-Type': 'application/json',
                                    },
                                    body: JSON.stringify({
                                        username: username,
                                        password: password
                                    })
                                })
                                .then(response => {
                                    if (response.ok) {
                                        const urlParams = new URLSearchParams(window.location.search);
                                        window.location.href = urlParams.get('redirect_url') || '/profile';
                                        return;
                                    }

                                    return response.json().then(data => {
                                        alert(data.message);
                                    });
                                })
                                .catch(error => {
                                    console.error('Error:', error);
                                    alert('Authentication failed. Please check your credentials and try again.');
                                });
                            });
                        });
                    </script>
                </div>
            </div>
        </div>
    </div>
</div><!-- end user -->
    {% endblock %}
