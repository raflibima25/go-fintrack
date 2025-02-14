<script>
export default {
  name: 'BaseAuth',
  props: {
    illustrationSrc: {
        type: String,
        required: true
    }
  },
}
</script>

<template>
    <div class="relative isolate min-h-screen bg-white">
      <!-- Gradient Background -->
      <div class="absolute inset-x-0 top-0 -z-10 transform-gpu overflow-hidden blur-3xl sm:top-[-20rem]">
        <div 
          class="relative left-1/2 -z-10 aspect-[1155/678] w-[36.125rem] -translate-x-1/2 bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] opacity-30 sm:w-[72.1875rem]" 
          style="clip-path: polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)"
        />
      </div>
  
      <!-- Content Container -->
      <div class="relative flex min-h-screen">
        <!-- Left side - illustration -->
        <div class="hidden lg:flex lg:w-1/2 items-center justify-end p-12">
          <div class="w-full max-w-xl flex items-center justify-end">
            <div class="w-full max-w-md space-y-4">
              <router-link to="/">
                <img src="../assets/logo-fintrack-new.webp" alt="Fintrack Logo" class="w-24 h-auto mx-auto -mb-20" />
                <img 
                  :src="illustrationSrc" 
                  alt="Authentication Illustration"
                  class="max-w-md"
                />
              </router-link>
            </div>
          </div>
        </div>
  
        <!-- Right side - auth form -->
        <div class="w-full lg:w-1/2 flex items-center justify-center lg:justify-start p-8">
          <div class="max-w-md w-full">
            <!-- Mobile view illustration -->
            <div class="lg:hidden text-center mb-8">
              <div class="w-full max-w-md space-y-4">
                <router-link to="/">
                  <img src="../assets/logo-fintrack-new.webp" alt="Fintrack Logo" class="w-24 h-auto mx-auto -mb-10" />
                  <img 
                    :src="illustrationSrc" 
                    alt="Authentication Illustration"
                    class="w-64 h-auto mx-auto"
                  />
              </router-link>
              </div>
            </div>
  
            <!-- Auth component slot -->
            <slot></slot>
          </div>
        </div>
      </div>
    </div>
</template>