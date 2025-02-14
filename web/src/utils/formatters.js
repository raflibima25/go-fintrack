export const formatCurrency = (value) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0
  }).format(value)
}

export const formatDate = (dateInput) => {
  try {
    const date = dateInput instanceof Date ? dateInput : new Date(dateInput)

    if (isNaN(date.getTime())) {
      throw new Error('Invalid date')
    }

    return date.toLocaleDateString('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
  } catch (err) {
    console.error('Error in formatDate:', err)
    return new Date().toLocaleDateString('en-US', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
  }
}

export const parseDate = (dateString) => {
  try {
    const date = new Date(dateString)
    return isNaN(date.getTime()) ? new Date() : date
  } catch {
    return new Date()
  }
}
