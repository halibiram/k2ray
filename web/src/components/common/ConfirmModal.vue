<template>
  <TransitionRoot appear :show="isOpen" as="template">
    <Dialog as="div" @close="closeModal" class="relative z-10">
      <TransitionChild
        as="template"
        enter="duration-300 ease-out"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="duration-200 ease-in"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <div class="fixed inset-0 bg-black bg-opacity-25" />
      </TransitionChild>

      <div class="fixed inset-0 overflow-y-auto">
        <div
          class="flex min-h-full items-center justify-center p-4 text-center"
        >
          <TransitionChild
            as="template"
            enter="duration-300 ease-out"
            enter-from="opacity-0 scale-95"
            enter-to="opacity-100 scale-100"
            leave="duration-200 ease-in"
            leave-from="opacity-100 scale-100"
            leave-to="opacity-0 scale-95"
          >
            <DialogPanel
              class="w-full max-w-md transform overflow-hidden rounded-2xl bg-white dark:bg-gray-800 p-6 text-left align-middle shadow-xl transition-all"
            >
              <DialogTitle
                as="h3"
                class="text-lg font-medium leading-6 text-gray-900 dark:text-gray-100"
              >
                {{ title }}
              </DialogTitle>
              <div class="mt-2">
                <p class="text-sm text-gray-500 dark:text-gray-400">
                  {{ message }}
                </p>
              </div>

              <div class="mt-4 flex justify-end space-x-2">
                <button
                  type="button"
                  class="inline-flex justify-center rounded-md border border-transparent bg-gray-200 dark:bg-gray-600 px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-100 hover:bg-gray-300 dark:hover:bg-gray-500 focus:outline-none"
                  @click="closeModal"
                >
                  {{ cancelButtonText }}
                </button>
                <button
                  type="button"
                  class="inline-flex justify-center rounded-md border border-transparent bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700 focus:outline-none"
                  @click="confirmAction"
                >
                  {{ confirmButtonText }}
                </button>
              </div>
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script setup lang="ts">
import {
  TransitionRoot,
  TransitionChild,
  Dialog,
  DialogPanel,
  DialogTitle,
} from '@headlessui/vue'

interface Props {
  isOpen: boolean
  title?: string
  message?: string
  confirmButtonText?: string
  cancelButtonText?: string
}

withDefaults(defineProps<Props>(), {
  isOpen: false,
  title: 'Confirm Action',
  message: 'Are you sure you want to proceed?',
  confirmButtonText: 'Confirm',
  cancelButtonText: 'Cancel',
})

const emit = defineEmits(['confirm', 'close'])

const closeModal = () => {
  emit('close')
}

const confirmAction = () => {
  emit('confirm')
}
</script>