<script lang="ts">
	import { onMount } from 'svelte';

	let scores: number[] = [];
	let currentPage = 1;
	const pageSize = 20;

	let paginatedScores: number[] = [];

	const apiUrl = 'https://your-api.com/scores'; // <-- Replace this

	async function fetchScores() {
		try {
			const response = await fetch(apiUrl);
			const data = await response.json();

			scores = data.map((item: { score: number }) => item.score).sort((a: any, b: any) => b - a);

			updatePage();
		} catch (error) {
			console.error('Error fetching scores:', error);
		}
	}

	function updatePage() {
		const start = (currentPage - 1) * pageSize;
		const end = start + pageSize;
		paginatedScores = scores.slice(start, end);
	}

	function nextPage() {
		if (currentPage * pageSize < scores.length) {
			currentPage++;
			updatePage();
		}
	}

	function prevPage() {
		if (currentPage > 1) {
			currentPage--;
			updatePage();
		}
	}

	onMount(fetchScores);
</script>

<div class="max-w-3xl mx-auto px-4 py-10">
	<h2 class="text-3xl font-bold text-center mb-6">🏆 Highscores</h2>

	<div class="overflow-x-auto rounded-lg shadow-lg bg-white">
		<table class="min-w-full text-sm text-left border-collapse">
			<thead class="bg-gray-100 text-gray-600 uppercase text-xs">
				<tr>
					<th class="px-6 py-4">Score</th>
				</tr>
			</thead>
			<tbody class="divide-y divide-gray-200">
				{#each paginatedScores as score}
					<tr class="hover:bg-gray-50 transition">
						<td class="px-6 py-4 font-mono text-gray-800">{score}</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>

	<div class="mt-6 flex justify-between items-center">
		<button
			on:click={prevPage}
			class="px-4 py-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300 disabled:opacity-50"
			disabled={currentPage === 1}
		>
			⬅ Previous
		</button>
		<span class="text-sm text-gray-600">Page {currentPage}</span>
		<button
			on:click={nextPage}
			class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
			disabled={currentPage * pageSize >= scores.length}
		>
			Next ➡
		</button>
	</div>
</div>
