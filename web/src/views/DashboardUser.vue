<script setup>
import { onMounted } from "vue";
import { useAuth } from "../composables/useAuth";
import UserLayout from "../layouts/UserLayout.vue";
import FinancialOverview from "../components/dashboard/FinancialOverview.vue";
import ExpenseAnalysis from "../components/dashboard/ExpenseAnalysis.vue";
import RecentActivity from "../components/dashboard/RecentActivity.vue";

const { checkAuth } = useAuth();

onMounted(() => {
  checkAuth('user');
});
</script>

<template>
  <user-layout>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- main container -->
      <div class="py-6">
        <!-- header section -->
        <div class="flex justify-between items-center mb-6">
          <h1 class="text-2xl font-bold">Dashboard</h1>
          <p class="text-gray-600">Welcome back, User!</p>
        </div>

        <!-- dashboard grid -->
        <div class="grid grid-cols-12 gap-6">
          <!-- financial overview -->
          <div class="col-span-12">
            <FinancialOverview />
          </div>

          <!-- expense analysis charts -->
          <div class="col-span-12 lg:col-span-8">
            <ExpenseAnalysis />
          </div>

          <!-- recent activity -->
          <div class="col-span-12 lg:col-span-4">
            <RecentActivity />
          </div>

          <!-- budget monitoring -->
          <div class="col-span-12 lg:col-span-8">
            <BudgetMonitoring />
          </div>

          <!-- financial goals -->
          <div class="col-span-12 lg:col-span-8">
            <FinancialGoals />
          </div>

          <!-- insights -->
          <div class="col-span-12 lg:col-span-8">
            <Insight />
          </div>
        </div>
      </div>
    </div>
  </user-layout>
</template>