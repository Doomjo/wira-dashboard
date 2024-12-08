<template>
  <div class="container mx-auto px-4 py-8">
    <!-- Search Bar -->
    <div class="mb-6">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search by username..."
        class="w-full px-4 py-2 border rounded-lg focus:outline-none focus:border-blue-500"
        @input="handleSearch"
      />
    </div>

    <!-- Players Table -->
    <div class="overflow-x-auto bg-white rounded-lg shadow">
      <table class="min-w-full table-auto">
        <thead class="bg-gray-100">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Rank</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Username</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Character Class</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Score</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="(player, index) in paginatedPlayers" :key="player.char_id" class="hover:bg-gray-50">
            <td class="px-6 py-4 whitespace-nowrap">{{ calculateRank(index) }}</td>
            <td class="px-6 py-4 whitespace-nowrap">{{ player.username }}</td>
            <td class="px-6 py-4 whitespace-nowrap">{{ player.class_name }}</td>
            <td class="px-6 py-4 whitespace-nowrap">{{ player.reward_score }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div class="mt-4 flex justify-between items-center">
      <div class="text-sm text-gray-700">
        Showing {{ paginationStart + 1 }} to {{ paginationEnd }} of {{ filteredPlayers.length }} entries
      </div>
      <div class="flex space-x-2">
        <button
          @click="currentPage--"
          :disabled="currentPage === 1"
          class="px-4 py-2 border rounded-lg disabled:opacity-50"
        >
          Previous
        </button>
        <button
          @click="currentPage++"
          :disabled="currentPage >= totalPages"
          class="px-4 py-2 border rounded-lg disabled:opacity-50"
        >
          Next
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'PlayerDashboard',
  data() {
    return {
      players: [],
      searchQuery: '',
      currentPage: 1,
      itemsPerPage: 10,
      loading: false,
      error: null
    };
  },
  computed: {
    filteredPlayers() {
      if (!this.searchQuery) return this.players;
      const query = this.searchQuery.toLowerCase();
      return this.players.filter(player => 
        player.username.toLowerCase().includes(query)
      );
    },
    paginatedPlayers() {
      const start = (this.currentPage - 1) * this.itemsPerPage;
      const end = start + this.itemsPerPage;
      return this.filteredPlayers.slice(start, end);
    },
    totalPages() {
      return Math.ceil(this.filteredPlayers.length / this.itemsPerPage);
    },
    paginationStart() {
      return (this.currentPage - 1) * this.itemsPerPage;
    },
    paginationEnd() {
      return Math.min(this.paginationStart + this.itemsPerPage, this.filteredPlayers.length);
    }
  },
  methods: {
    async fetchPlayers() {
      try {
        this.loading = true;
        const response = await axios.get('http://localhost:8080/api/players');
        this.players = response.data.sort((a, b) => b.reward_score - a.reward_score);
      } catch (err) {
        this.error = 'Error fetching players data';
        console.error(err);
      } finally {
        this.loading = false;
      }
    },
    handleSearch() {
      this.currentPage = 1; // Reset to first page when searching
    },
    calculateRank(index) {
      return this.paginationStart + index + 1;
    }
  },
  created() {
    this.fetchPlayers();
  }
};
</script>

<style scoped>
@media (max-width: 640px) {
  table {
    display: block;
    overflow-x: auto;
    white-space: nowrap;
  }
  
  .container {
    padding: 1rem;
  }
  
  th, td {
    padding: 0.5rem;
  }
}
</style>
