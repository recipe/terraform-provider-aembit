provider "aembit" {
}

resource "aembit_client_workload" "test" {
    name = "Unit Test 1"
    description = "Acceptance Test client workload"
    is_active = true
    identities = [
        {
            type = "k8sNamespace"
            value = "unittest1namespace"
        },
    ]
}

