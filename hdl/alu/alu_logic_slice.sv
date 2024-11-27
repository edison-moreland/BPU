//visualization: netlistsvg

module alu_logic_slice #(
    parameter N = 1
) (
    input  logic [N-1:0] a, b,
    input  logic         enable_xor, enable_or, enable_and, enable_not,
    output logic [N-1:0] result,
);

    logic [N-1:0] result_xor, result_or, result_and, result_not;

    always_comb begin
        result_xor <= (a ^ b);
        result_and <= (a & b);
        result_not <= (~a);
        result_or  <= (a | b);

        result <= ({N{enable_xor}} & result_xor) |
                  ({N{enable_or}}  & result_or ) |
                  ({N{enable_and}} & result_and) |
                  ({N{enable_not}} & result_not);
    end

endmodule