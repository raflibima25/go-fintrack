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
    <div class="min-h-screen bg-white flex">
        <!-- left side - illustration -->
        <div class="hidden lg:flex lg:w-1/2 items-center justify-end bg-white p-12">
            <div class="w-full max-w-xl flex items-center justify-end">
                <img 
                    :src="illustrationSrc" 
                    alt="Authentication Illustration"
                    class="max-w-md"
                />
            </div>
        </div>

        <!-- right side - auth form -->
        <div class="w-full lg:w-1/2 flex items-center justify-start p-8">
            <div class="max-w-md w-full">
                <!-- logo mobile view -->
                <div class="lg:hidden text-center mb-8">
                    <img 
                        :src="illustrationSrc" 
                        alt="Authentication Illustration"
                        class="w-64 h-auto mx-auto"
                    />
                </div>

                <!-- slot komponen auth -->
                <slot></slot>
            </div>
        </div>
    </div>
</template>