<!-- LandingPage.vue -->
<script setup>
import { ref } from 'vue'
import { 
  PresentationChartBarIcon,
  CalendarDateRangeIcon,
  LockClosedIcon,
  ArrowPathIcon,
  Bars3Icon,
  XMarkIcon 
} from '@heroicons/vue/24/outline'
// logo fintrack
import logoFintrack from '../assets/logo-fintrack-new.webp'
// import blog images
import aionyHaust from '../assets/blog/blog-avatar/aiony-haust-3TLl_97HNJo-unsplash.webp'
import ianDooley from '../assets/blog/blog-avatar/ian-dooley-d1UPkiFd04A-unsplash.webp'
import michaelDam from '../assets/blog/blog-avatar/michael-dam-mEZ3PoFGs_k-unsplash.webp'
import mihaiMoisa from '../assets/blog/mihai-moisa-Djtc1T38-GY-unsplash.webp'
import nathanDumlaoA0t from '../assets/blog/nathan-dumlao-A0tNoiSq4mo-unsplash.webp'
import nathanDumlaoY3A from '../assets/blog/nathan-dumlao-Y3AqmbmtLQI-unsplash.webp'

const mobileMenuOpen = ref(false)

// Navigation data
const navigation = [
  { name: 'Product', href: '#' },
  { name: 'Features', href: '#' },
  { name: 'Marketplace', href: '#' },
  { name: 'Company', href: '#' },
]

// Features data
const features = [
  {
    name: 'Smart Budgeting',
    description: 'Track your income and expenses effortlessly with automated categorization and insightful reports.',
    icon: PresentationChartBarIcon,
  },
  {
    name: 'Secure Transactions',
    description: 'Your financial data is encrypted and protected with the latest security standards for peace of mind.',
    icon: LockClosedIcon,
  },
  {
    name: 'Automated Savings',
    description: 'Set savings goals and let our system automatically allocate funds to help you achive them faster.',
    icon: ArrowPathIcon,
  },
  {
    name: 'Financial Planning',
    description: 'Plan ahead with expense forecasts, bill reminders, and investment tracking for long-term stability.',
    icon: CalendarDateRangeIcon,
  },
]

// Blog data
const blogs = [
  {
    title: 'Boost your conversion rate',
    date: 'Jan 16, 2025',
    author: {
      name: 'Michael Foster',
      imageUrl: aionyHaust
    },
    imageUrl: mihaiMoisa,
  },
  {
    title: 'How to use search engine optimization to drive sales',
    date: 'Feb 10, 2025', 
    author: {
      name: 'Lindsay Walton',
      imageUrl: ianDooley
    },
    imageUrl: nathanDumlaoA0t,
  },
  {
    title: 'Improve your customer experience',
    date: 'Feb 12, 2025',
    author: {
      name: 'Tom Cook',
      imageUrl: michaelDam
    },
    imageUrl: nathanDumlaoY3A,
  }
]

// Footer navigation
const footerNavigation = {
  solutions: [
    { name: 'Marketing', href: '#' },
    { name: 'Analytics', href: '#' },
    { name: 'Commerce', href: '#' },
    { name: 'Insights', href: '#' },
  ],
  support: [
    { name: 'Pricing', href: '#' },
    { name: 'Documentation', href: '#' },
    { name: 'Guides', href: '#' },
    { name: 'API Status', href: '#' },
  ],
  company: [
    { name: 'About', href: '#' },
    { name: 'Blog', href: '#' },
    { name: 'Jobs', href: '#' },
    { name: 'Press', href: '#' },
    { name: 'Partners', href: '#' },
  ],
  legal: [
    { name: 'Claim', href: '#' },
    { name: 'Privacy', href: '#' },
    { name: 'Terms', href: '#' },
  ],
}

// FAQ data
const faqs = [
  {
    question: "What is financial management?",
    answer: "Financial management refers to the strategic planning, organizing, directing, and controlling of financial resources to achieve financial goals efficiently."
  },
  {
    question: "Why is budgeting important?",
    answer: "Budgeting helps track income and expenses, ensuring you live within your means, save for the future, and avoid unnecessary debt."
  },
  {
    question: "How can I improve my financial health?",
    answer: "You can improve your financial health by creating a budget, saving regularly, managing debt wisely, and investing for long-term growth."
  },
  {
    question: "What are the benefits of automated savings?",
    answer: "Automated savings help you consistently set aside money for emergencies, retirement, and financial goals without having to think about it."
  },
  {
    question: "Is my financial data secure?",
    answer: "Yes, our platform uses advanced encryption and security protocols to ensure your financial data is safe and protected."
  }
];

// Get the current year
const currentYear = new Date().getFullYear();
</script>

<template>
  <div class="bg-white">
    <!-- Header -->
    <header class="absolute inset-x-0 top-0 z-50">
      <nav class="flex items-center justify-between p-6 lg:px-8" aria-label="Global">
        <!-- Logo -->
        <div class="flex lg:flex-1">
          <a href="#" class="-m-1.5 p-1.5">
            <span class="sr-only">Fintrack</span>
            <img class="h-10 w-auto" :src=logoFintrack alt="Fintrack Logo" />
          </a>
        </div>

        <!-- Mobile menu button -->
        <div class="flex lg:hidden">
          <button 
            type="button" 
            class="-m-2.5 inline-flex items-center justify-center rounded-md p-2.5 text-gray-700"
            @click="mobileMenuOpen = true"
          >
            <span class="sr-only">Open main menu</span>
            <Bars3Icon class="h-6 w-6" aria-hidden="true" />
          </button>
        </div>

        <!-- Desktop navigation -->
        <div class="hidden lg:flex lg:gap-x-12">
          <a v-for="item in navigation" 
             :key="item.name" 
             :href="item.href" 
             class="text-sm font-semibold leading-6 text-gray-900"
          >
            {{ item.name }}
          </a>
        </div>

        <!-- Desktop login button -->
        <div class="hidden lg:flex lg:flex-1 lg:justify-end">
          <a href="/login" class="text-sm font-semibold leading-6 text-gray-900">
            Log in <span aria-hidden="true">→</span>
          </a>
        </div>
      </nav>

      <!-- Mobile menu -->
      <div v-if="mobileMenuOpen" class="lg:hidden">
        <div class="fixed inset-0 z-50"></div>
        <div class="fixed inset-y-0 right-0 z-50 w-full overflow-y-auto bg-white px-6 py-6 sm:max-w-sm sm:ring-1 sm:ring-gray-900/10">
          <div class="flex items-center justify-between">
            <a href="#" class="-m-1.5 p-1.5">
              <span class="sr-only">Fintrack</span>
              <img class="h-10 w-auto" :src=logoFintrack alt="Fintrack Logo" />
            </a>
            <button 
              type="button" 
              class="-m-2.5 rounded-md p-2.5 text-gray-700"
              @click="mobileMenuOpen = false"
            >
              <span class="sr-only">Close menu</span>
              <XMarkIcon class="h-6 w-6" aria-hidden="true" />
            </button>
          </div>
          <div class="mt-6 flow-root">
            <div class="-my-6 divide-y divide-gray-500/10">
              <div class="space-y-2 py-6">
                <a v-for="item in navigation" 
                   :key="item.name" 
                   :href="item.href" 
                   class="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50"
                >
                  {{ item.name }}
                </a>
              </div>
              <div class="py-6">
                <a href="/login" class="-mx-3 block rounded-lg px-3 py-2.5 text-base font-semibold leading-7 text-gray-900 hover:bg-gray-50">
                  Log in
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </header>

    <!-- Hero section -->
    <div class="relative isolate pt-14">
      <div class="absolute inset-x-0 -top-40 -z-10 transform-gpu overflow-hidden blur-3xl sm:-top-80">
        <div 
          class="relative left-[calc(50%-11rem)] aspect-[1155/678] w-[36.125rem] -translate-x-1/2 rotate-[30deg] bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] opacity-30 sm:left-[calc(50%-30rem)] sm:w-[72.1875rem]" 
          style="clip-path: polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)"
        />
      </div>
      <div class="py-24 sm:py-32 lg:pb-40">
        <div class="mx-auto max-w-7xl px-6 lg:px-8">
          <div class="mx-auto max-w-2xl text-center">
            <h1 class="text-4xl font-bold tracking-tight text-gray-900 sm:text-6xl">
              Manage your finances smartly
            </h1>
            <p class="mt-6 text-lg leading-8 text-gray-600">
              Manage your expenses, increase your savings, and achieve your financial goals with ease. Safe, efficient and practical financial solutions for a more peaceful life.
            </p>
            <div class="mt-10 flex items-center justify-center gap-x-6">
              <a href="/login" class="rounded-md bg-indigo-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">
                Get started
              </a>
              <a href="#" class="text-sm font-semibold leading-6 text-gray-900">
                Learn more <span aria-hidden="true">→</span>
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Features section -->
    <div class="mx-auto max-w-7xl px-6 sm:mt-24 lg:px-8">
      <div class="mx-auto max-w-2xl text-center">
        <h2 class="text-base font-semibold leading-7 text-indigo-600">Tool faster</h2>
        <p class="mt-2 text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
          Everything you need to manage your finances with ease
        </p>
        <p class="mt-6 text-lg leading-8 text-gray-600">
          Take control of your financial future with powerful tools designed to help you budget, save, and invest wisely.
        </p>
      </div>
      <div class="mx-auto mt-16 max-w-2xl sm:mt-20 lg:mt-24 lg:max-w-none">
        <dl class="grid max-w-xl grid-cols-1 gap-x-8 gap-y-16 lg:max-w-none lg:grid-cols-4">
          <div v-for="feature in features" :key="feature.name" class="flex flex-col">
            <dt class="flex items-center gap-x-3 text-base font-semibold leading-7 text-gray-900">
              <component :is="feature.icon" class="h-5 w-5 flex-none text-indigo-600" aria-hidden="true" />
              {{ feature.name }}
            </dt>
            <dd class="mt-4 flex flex-auto flex-col text-base leading-7 text-gray-600">
              <p class="flex-auto">{{ feature.description }}</p>
            </dd>
          </div>
        </dl>
      </div>
    </div>

    <!-- Blog section -->
    <div class="py-24 sm:pt-48">
      <div class="mx-auto max-w-7xl px-6 lg:px-8">
        <div class="mx-auto max-w-4xl text-center">
          <h2 class="text-base font-semibold leading-7 text-indigo-600">Blog</h2>
          <p class="mt-2 text-4xl font-bold tracking-tight text-gray-900 sm:text-5xl">
            From the blog
          </p>
          <p class="mt-6 text-lg leading-8 text-gray-600">
            Learn how to grow your business with our expert advice.
          </p>
        </div>

        <div class="mx-auto mt-16 grid max-w-2xl grid-cols-1 gap-x-8 gap-y-20 lg:mx-0 lg:max-w-none lg:grid-cols-3">
          <article v-for="post in blogs" :key="post.title" class="flex flex-col items-start relative isolate overflow-hidden rounded-2xl h-[400px]">
            <!-- Background image -->
            <img
              :src="post.imageUrl"
              :alt="post.title"
              class="absolute inset-0 h-full w-full object-cover"
            />
            <!-- Gradient overlay -->
            <div class="absolute inset-0 bg-gradient-to-t from-gray-900/75 to-gray-900/0"></div>
            
            <div class="relative p-8 flex flex-col h-full justify-end">
              <!-- Date and author -->
              <div class="flex items-center gap-x-4 text-sm leading-6 text-white">
                <time :datetime="post.date">{{ post.date }}</time>
                <div class="flex items-center gap-x-2">
                  <img :src="post.author.imageUrl" :alt="post.author.name" class="h-6 w-6 rounded-full bg-white/10" />
                  <span>{{ post.author.name }}</span>
                </div>
              </div>
              <!-- Title -->
              <h3 class="mt-3 text-lg font-semibold leading-6 text-white">
                <a href="#">
                  <span class="absolute inset-0"></span>
                  {{ post.title }}
                </a>
              </h3>
            </div>
          </article>
        </div>
      </div>
    </div>

    <!-- FAQ section -->
    <div class="mx-auto max-w-2xl divide-y divide-gray-900/10 px-6 pb-8 sm:pb-24 sm:pt-12 lg:max-w-7xl lg:px-8 lg:pb-32">
      <h2 class="text-2xl font-bold leading-10 tracking-tight text-gray-900">Frequently asked questions</h2>
      <dl class="mt-10 space-y-8 divide-y divide-gray-900/10">
        <div v-for="faq in faqs" :key="faq.question" class="pt-8 lg:grid lg:grid-cols-12 lg:gap-8">
          <dt class="text-base font-semibold leading-7 text-gray-900 lg:col-span-5">{{ faq.question }}</dt>
          <dd class="mt-4 lg:col-span-7 lg:mt-0">
            <p class="text-base leading-7 text-gray-600">{{ faq.answer }}</p>
          </dd>
        </div>
      </dl>
    </div>

    <!-- CTA section -->
    <div class="relative isolate px-6 py-32 sm:mt-24 sm:py-40 lg:px-8">
      <div class="absolute inset-x-0 top-0 -z-10 transform-gpu overflow-hidden blur-3xl sm:top-[-20rem]">
        <div 
          class="relative left-1/2 -z-10 aspect-[1155/678] w-[36.125rem] -translate-x-1/2 bg-gradient-to-tr from-[#ff80b5] to-[#9089fc] opacity-30 sm:w-[72.1875rem]" 
          style="clip-path: polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)"
        />
      </div>
      <div class="mx-auto max-w-2xl text-center">
        <h2 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
          Take Control of Your Finances.<br>Start Managing Smarter Today.
        </h2>
        <p class="mx-auto mt-6 max-w-xl text-lg leading-8 text-gray-600">
          Simplify your financial management, track expenses, set budgets, and achieve your financial goals effortlessly.
        </p>
        <div class="mt-10 flex items-center justify-center gap-x-6">
          <a href="/login" class="rounded-md bg-indigo-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">
            Get started
          </a>
          <a href="#" class="text-sm font-semibold leading-6 text-gray-900">
            Learn more <span aria-hidden="true">→</span>
          </a>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <footer class="bg-white" aria-labelledby="footer-heading">
      <h2 id="footer-heading" class="sr-only">Footer</h2>
      <div class="mx-auto max-w-7xl px-6 pb-8 pt-16 sm:pt-24 lg:px-8 lg:pt-32">
        <div class="xl:grid xl:grid-cols-3 xl:gap-8">
          <div class="space-y-8">
            <img class="h-7" :src=logoFintrack alt="Fintrack" />
            <p class="text-sm leading-6 text-gray-600">
              Making the world a better place through constructing elegant hierarchies.
            </p>
            <div class="flex space-x-6">
              <a v-for="item in footerNavigation.social" :key="item.name" :href="item.href" class="text-gray-400 hover:text-gray-500">
                <span class="sr-only">{{ item.name }}</span>
                <component :is="item.icon" class="h-6 w-6" aria-hidden="true" />
              </a>
            </div>
          </div>
          <div class="mt-16 grid grid-cols-2 gap-8 xl:col-span-2 xl:mt-0">
            <div class="md:grid md:grid-cols-2 md:gap-8">
              <div>
                <h3 class="text-sm font-semibold leading-6 text-gray-900">Solutions</h3>
                <ul role="list" class="mt-6 space-y-4">
                  <li v-for="item in footerNavigation.solutions" :key="item.name">
                    <a :href="item.href" class="text-sm leading-6 text-gray-600 hover:text-gray-900">
                      {{ item.name }}
                    </a>
                  </li>
                </ul>
              </div>
              <div class="mt-10 md:mt-0">
                <h3 class="text-sm font-semibold leading-6 text-gray-900">Support</h3>
                <ul role="list" class="mt-6 space-y-4">
                  <li v-for="item in footerNavigation.support" :key="item.name">
                    <a :href="item.href" class="text-sm leading-6 text-gray-600 hover:text-gray-900">
                      {{ item.name }}
                    </a>
                  </li>
                </ul>
              </div>
            </div>
            <div class="md:grid md:grid-cols-2 md:gap-8">
              <div>
                <h3 class="text-sm font-semibold leading-6 text-gray-900">Company</h3>
                <ul role="list" class="mt-6 space-y-4">
                  <li v-for="item in footerNavigation.company" :key="item.name">
                    <a :href="item.href" class="text-sm leading-6 text-gray-600 hover:text-gray-900">
                      {{ item.name }}
                    </a>
                  </li>
                </ul>
              </div>
              <div class="mt-10 md:mt-0">
                <h3 class="text-sm font-semibold leading-6 text-gray-900">Legal</h3>
                <ul role="list" class="mt-6 space-y-4">
                  <li v-for="item in footerNavigation.legal" :key="item.name">
                    <a :href="item.href" class="text-sm leading-6 text-gray-600 hover:text-gray-900">
                      {{ item.name }}
                    </a>
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
        <div class="mt-16 border-t border-gray-900/10 pt-8 sm:mt-20 lg:mt-24">
          <p class="text-xs leading-5 text-gray-500">
            &copy; <span>{{ currentYear }}</span> Fintrack, Inc. All rights reserved.
          </p>
        </div>
      </div>
    </footer>

  </div>
</template>